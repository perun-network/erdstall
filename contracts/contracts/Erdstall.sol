// SPDX-License-Identifier: Apache-2.0

pragma solidity ^0.7.0;
pragma experimental ABIEncoderV2;

import "./Sig.sol";

contract Erdstall {
    // The epoch-balance statements signed by the TEE.
    struct Balance {
        uint64 epoch;
        address account;
        uint256 value;
    }

    uint64 constant notFrozen = uint64(-2); // use 2nd-highest number to indicate not-frozen

    // Parameters set during deployment.
    address public immutable tee; // yummi 🍵
    uint64 public immutable bigBang; // start of first epoch
    uint64 public immutable phaseDuration; // number of blocks of one epoch phase
    uint64 public immutable responseDuration; // operator response grace period

    mapping(uint64 => mapping(address => uint256)) public deposits; // epoch => account => balance value
    mapping(uint64 => mapping(address => uint256)) public exits; // epoch => account => balance value
    mapping(uint64 => mapping(address => uint256)) public challenges; // epoch => account => recovery value
    mapping(uint64 => uint256) public numChallenges; // epoch => numChallenges
    mapping(address => bool) public frozenWithdrawals; // account => withdrawn-flag
    uint64 public frozenEpoch = notFrozen; // epoch at which contract was frozen

    event Deposited(uint64 indexed epoch, address indexed account, uint256 value);
    event Exiting(uint64 indexed epoch, address indexed account, uint256 value);
    event Withdrawn(uint64 indexed epoch, address indexed account, uint256 value);
    event Challenged(uint64 indexed epoch, address indexed account);
    event Frozen(uint64 indexed epoch);

    constructor(address _tee, uint64 _phaseDuration, uint64 _responseDuration) {
        // responseDuration should be at most half the phaseDuration
        require(2 * _responseDuration <= _phaseDuration, "responseDuration too long");
        tee = _tee;
        bigBang = uint64(block.number);
        phaseDuration = _phaseDuration;
        responseDuration = _responseDuration;
    }

    modifier onlyAlive {
        require(!isFrozen(), "plasma frozen");
        // in case freeze wasn't called yet...
        require(!isLastEpochChallenged(), "plasma freezing");
        _;
    }

    //
    // Normal Operation
    //

    function deposit() external payable onlyAlive {
        uint64 epoch = depositEpoch();
        deposits[epoch][msg.sender] += msg.value;

        emit Deposited(epoch, msg.sender, msg.value);
    }

    // exit lets a user exit and the end of the epoch's exit period.
    // sig must be signature created by signText(keccak256(abi.encode(balance))).
    // For now, only full exits are allowed.
    //
    // exit is also used to answer challenges.
    function exit(Balance calldata balance, bytes calldata sig) external onlyAlive {
        require(balance.epoch == exitEpoch(), "exit: wrong epoch");
        verifyBalance(balance, sig);
        if (challenges[balance.epoch][balance.account] == 0) {
            // if not challenged, only user can exit
            require(balance.account == msg.sender, "exit: wrong sender");
        } else {
            // reset challenge if this is a challenge response
            challenges[balance.epoch][balance.account] = 0;
            numChallenges[balance.epoch]--;
        }

        exits[balance.epoch][balance.account] = balance.value;

        emit Exiting(balance.epoch, balance.account, balance.value);
    }

    function withdraw(uint64 epoch) external onlyAlive {
        // can only withdraw after exit period
        require(epoch < exitEpoch(), "withdraw: too early");

        uint256 value = exits[epoch][msg.sender];
        require(value > 0, "nothing left to withdraw");
        exits[epoch][msg.sender] = 0;

        msg.sender.transfer(value);
        emit Withdrawn(epoch, msg.sender, value);
    }

    //
    // Challenge Functions
    //

    // Challenges the operator to post the current exit epoch's balance statement.
    // The user needs to pass the latest balance proof, that is, of the just
    // sealed epoch, to proof that they are part of the system.
    //
    // After a challenge is opened, the operator (anyone, actually) can respond
    // to the challenge using function `exit`.
    function challenge(Balance calldata balance, bytes calldata sig) external onlyAlive {
        require(balance.account == msg.sender, "challenge: wrong sender");
        require(balance.epoch == sealedEpoch(), "challenge: wrong epoch");
        verifyBalance(balance, sig);

        registerChallenge(balance.value);
    }

    // challengeDeposit should be called by a user if they deposited but never
    // received a deposit or balance proof from the operator.
    //
    // After a challenge is opened, the operator (anyone, actually) can respond
    // to the challenge using function `exit`.
    function challengeDeposit() external onlyAlive {
        registerChallenge(0);
    }

    function registerChallenge(uint256 recoveryBalance) internal {
        require(!isChallengeResponsePhase(), "in challenge response phase");
        uint64 epoch = exitEpoch();
        require(challenges[epoch][msg.sender] == 0, "already challenged");

        uint256 value = recoveryBalance + deposits[epoch][msg.sender];
        require(value > 0, "no value in system");

        challenges[epoch][msg.sender] = value;
        numChallenges[epoch]++;

        emit Challenged(epoch, msg.sender);
    }

    // withdrawChallenge lets open challengers withdraw all funds locked in the
    // frozen contract. The funds were already determined when the challenge was
    // posted using either `challenge` or `challengeDeposit`.
    //
    // Implicitly calls ensureFrozen to ensure that the contract state is set to
    // frozen if the last epoch has an unanswered challenge.
    function withdrawChallenge() external {
        ensureFrozen();

        uint256 value = challenges[frozenEpoch+1][msg.sender];
        require(value > 0, "nothing left to withdraw (frozen)");

        _withdrawFrozen(value);
    }

    // withdrawFrozen lets non-challengers withdraw all funds locked in the
    // frozen contract. Parameter `balance` needs to be the balance proof of the
    // last unchallenged epoch.
    //
    // Implicitly calls ensureFrozen to ensure that the contract state is set to
    // frozen if the last epoch has an unanswered challenge.
    function withdrawFrozen(Balance calldata balance, bytes calldata sig) external {
        ensureFrozen();

        require(balance.account == msg.sender, "withdrawFrozen: wrong sender");
        require(balance.epoch == frozenEpoch, "withdrawFrozen: wrong epoch");
        verifyBalance(balance, sig);

        // Also recover deposits from broken epoch
        uint256 value = balance.value + deposits[frozenEpoch+1][msg.sender];

        _withdrawFrozen(value);
    }

    function _withdrawFrozen(uint256 value) internal {
        require(!frozenWithdrawals[msg.sender], "already withdrawn (frozen)");
        frozenWithdrawals[msg.sender] = true;

        msg.sender.transfer(value);
        emit Withdrawn(frozenEpoch, msg.sender, value);
    }

    // ensureFrozen ensures that the state of the contract is set to frozen if
    // the last epoch has at least one unanswered challenge.
    //
    // It is implicitly called by withdrawFrozen but can be called seperately if
    // the contract should be frozen before anyone wants to withdraw.
    function ensureFrozen() public {
        if (isFrozen()) { return; }
        require(isLastEpochChallenged(), "no challenge in last epoch");

        // freezing to previous epoch
        uint64 epoch = freezingEpoch() - 1;
        frozenEpoch = epoch;

        emit Frozen(epoch);
    }

    function isLastEpochChallenged() internal view returns (bool) {
        return numChallenges[freezingEpoch()] > 0;
    }

    function isFrozen() internal view returns (bool) {
        return frozenEpoch != notFrozen;
    }

    function isChallengeResponsePhase() internal view returns (bool) {
        // the last responseDuration blocks of each exit phase are reserved for
        // challenge responses.
        return phaseDuration - ((uint64(block.number) - bigBang) % phaseDuration) <= responseDuration;
    }

    //
    // Epoch Counter Abstractions
    //

    function depositEpoch() internal view returns (uint64) {
        return epoch();
    }

    function exitEpoch() internal view returns (uint64) {
        return epoch()-2;
    }

    function freezingEpoch() internal view returns (uint64) {
        return epoch()-3;
    }

    function sealedEpoch() internal view returns (uint64) {
        return epoch()-3;
    }

    // epoch returns the current epoch. It should not be used directly in public
    // functions, but the fooEoch functions instead, as they account for the
    // correct shifts.
    function epoch() internal view returns (uint64) {
        return (uint64(block.number) - bigBang) / phaseDuration;
    }

    function verifyBalance(Balance memory balance, bytes memory sig) public view {
        require(Sig.verify(encodeBalanceProof(balance), sig, tee), "invalid signature");
    }

    function encodeBalanceProof(Balance memory balance) public view returns (bytes memory) {
        return abi.encode(
            "ErdstallBalance",
            address(this),
            balance.epoch,
            balance.account,
            balance.value);
    }
}
