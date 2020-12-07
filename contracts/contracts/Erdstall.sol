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
    mapping(uint64 => mapping(address => bool)) public challenges; // epoch => account => challenge-flag
    mapping(uint64 => uint256) public numChallenges; // epoch => numChallenges
    uint64 public frozenEpoch = notFrozen; // epoch at which contract was frozen
    mapping(uint64 => mapping(address => bool)) public frozenWithdraws; // epoch => account => frozen-withdrawn-flag

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
    function exit(Balance calldata balance, bytes calldata sig) external onlyAlive {
        require(balance.epoch == exitEpoch());
        verifyBalance(balance, sig);
        if (!challenges[balance.epoch][balance.account]) {
            // if not challenged, only user can exit
            require(balance.account == msg.sender, "exit: wrong sender");
        } else {
            // reset challenge if this is a challenge response
            challenges[balance.epoch][balance.account] = false;
            numChallenges[balance.epoch]--;
        }

        exits[balance.epoch][balance.account] = balance.value;

        emit Exiting(balance.epoch, balance.account, balance.value);
    }

    function withdraw(uint64 epoch) external onlyAlive {
        // can only withdraw after exit period
        require(epoch < exitEpoch(), "withdraw: too early");

        uint256 value = exits[epoch][msg.sender];
        exits[epoch][msg.sender] = 0;

        msg.sender.transfer(value);
        emit Withdrawn(epoch, msg.sender, value);
    }

    //
    // Challenge Functions
    //

    // Challenges the operator to post the epoch's balance statement.
    // After a challenge is opened, the operator (anyone, actually) can respond
    // to the challenge using function exit.
    function challenge() external onlyAlive {
        require(!isChallengeResponsePhase(), "in challenge response phase");
        uint64 epoch = exitEpoch();
        require(!challenges[epoch][msg.sender], "already challenged");

        challenges[epoch][msg.sender] = true;
        numChallenges[epoch]++;

        emit Challenged(epoch, msg.sender);
    }

    function withdrawFrozen(Balance calldata balance, bytes calldata sig) external {
        verifyBalance(balance, sig);

        _withdrawFrozen(balance.value);
    }

    // recoverDeposit should be called by a user if they deposited but never
    // received a deposit or balance proof from the operator. They must have
    // launched an unanswered challenge in the epoch's exit phase so that
    // the contract can be frozen.
    function recoverDeposit() external {
        _withdrawFrozen(0);
    }

    function _withdrawFrozen(uint256 _value) internal {
        ensureFrozen();
        require(!frozenWithdraws[frozenEpoch][msg.sender], "already withdrawn (frozen)");

        uint256 value = _value + frozenDeposit();
        frozenWithdraws[frozenEpoch][msg.sender] = true;

        msg.sender.transfer(value);
        emit Withdrawn(frozenEpoch, msg.sender, value);
    }

    // ensureFrozen ensures that the state of the contract is set to frozen if
    // the last epoch has at least one unanswered challenge.
    function ensureFrozen() internal {
        if (isFrozen()) { return; }
        require(isLastEpochChallenged(), "no challenge in last epoch");

        // freezing to previous epoch
        uint64 epoch = freezingEpoch() - 1;
        frozenEpoch = epoch;

        emit Frozen(epoch);
    }

    function frozenDeposit() internal view returns (uint256) {
        return deposits[frozenEpoch+1][msg.sender];
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
