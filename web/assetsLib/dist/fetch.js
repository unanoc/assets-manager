"use strict";
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.httpPostFromBrowser = exports.fetchUniversal = void 0;
const axios_1 = __importDefault(require("axios"));
function errorHead(error) {
    return error.toString().substring(0, 400 - 1);
}
// Fetch a URL, either:
// - using axios, when used from 'backend'
// - using javascript fetch, when used from within a browser
function fetchUniversal(url, fromBrowser) {
    return __awaiter(this, void 0, void 0, function* () {
        try {
            if (fromBrowser) {
                return yield fetchBrowser(url);
            }
            return yield fetchBackend(url);
        }
        catch (ex) {
            console.log('Exception:', ex);
            return [500, errorHead(ex)];
        }
    });
}
exports.fetchUniversal = fetchUniversal;
function fetchBackend(url) {
    return __awaiter(this, void 0, void 0, function* () {
        const resp = yield axios_1.default.get(url);
        if (resp.status != 200) {
            console.log("ERROR: Non-OK status", resp.status, resp.statusText, url);
            return [resp.status, ""];
        }
        let text = "";
        try {
            text = JSON.stringify(resp.data);
        }
        catch (error) {
            text = resp.data;
        }
        return [resp.status, text];
    });
}
function fetchBrowser(url) {
    return __awaiter(this, void 0, void 0, function* () {
        let resp = yield fetch(url);
        if (resp.status != 200) {
            console.log("ERROR: Non-OK status", resp.status, resp.statusText, url);
            return [resp.status, ""];
        }
        const text = yield resp.text();
        return [resp.status, text];
    });
}
function httpPostFromBrowser(url, data) {
    return __awaiter(this, void 0, void 0, function* () {
        const options = {
            method: 'POST',
            body: JSON.stringify(data),
            headers: {
                'Content-Type': 'application/json'
            }
        };
        //console.log('options.body', options.body);
        const resp = yield fetch(url, options);
        if (resp.status != 200) {
            console.log("ERROR: Non-OK status", resp.status, resp.statusText, url);
            return [resp.status, ""];
        }
        const json = yield resp.json();
        return [resp.status, json];
    });
}
exports.httpPostFromBrowser = httpPostFromBrowser;
