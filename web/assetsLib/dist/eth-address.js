"use strict";
//import { reverseCase } from "./types"; // TODO
Object.defineProperty(exports, "__esModule", { value: true });
exports.isEthereumAddress = exports.isChecksum = exports.toChecksum = exports.toChecksumEthereum = exports.isChecksumEthereum = void 0;
const ethereum_checksum_address_1 = require("ethereum-checksum-address");
exports.isChecksumEthereum = (address) => ethereum_checksum_address_1.checkAddressChecksum(address);
exports.toChecksumEthereum = (address) => ethereum_checksum_address_1.toChecksumAddress(address);
function toChecksum(address, chain = "ethereum") {
    try {
        const checksumEthereum = exports.toChecksumEthereum(address);
        /*
        // special handling for Wanchain
        if (chain.toLowerCase() === "wanchain") {
            const checksumWanchain = reverseCase(checksumEthereum).replace("X", "x");
            return checksumWanchain;
        }
        */
        return checksumEthereum;
    }
    catch (error) {
        console.log('Exception:', error);
    }
    return address; // falbback
}
exports.toChecksum = toChecksum;
function isChecksum(address, chain = "ethereum") {
    /*
    // special handling for Wanchain
    if (chain.toLowerCase() === "wanchain") {
        const addressEthereum = reverseCase(address).replace("X", "x");
        return isChecksumEthereum(addressEthereum);
    }
    */
    return exports.isChecksumEthereum(address);
}
exports.isChecksum = isChecksum;
function isEthereumAddress(address) {
    if (!(address.length == 40 || (address.length == 42 && address.substring(0, 2) === '0x'))) {
        return false;
    }
    const check1 = toChecksum(address);
    if (toChecksum(check1) !== check1) {
        return false;
    }
    return true;
}
exports.isEthereumAddress = isEthereumAddress;
