// SPDX-License-Identifier: SEE LICENSE IN LICENSE
pragma solidity ^0.8.0;

contract Counter {
    uint256 public count;

    event CountIncrement(address indexed caller, uint256 newValue);
    event CountDecrement(address indexed caller, uint256 newValue);
    event CountReset(address indexed caller, uint256 newValue);

    function increment() external {
        count++;
        emit CountIncrement(msg.sender, count);
    }

    function decrement() external {
        require(count > 0, "Counter: count cannot be negative");
        count--;
        emit CountDecrement(msg.sender, count);
    }

    function getCount() external view returns (uint256) {
        return count;
    }

    function reset() external {
        count = 0;
        emit CountReset(msg.sender, count);
    }

}