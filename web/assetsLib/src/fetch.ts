import axios from "axios";


function errorHead(error): string {
    return error.toString().substring(0, 400-1);
}

// Fetch a URL, either:
// - using axios, when used from 'backend'
// - using javascript fetch, when used from within a browser
export async function fetchUniversal(url: string, fromBrowser: boolean): Promise<[number, string]> {
    try {
        if (fromBrowser) {
            return await fetchBrowser(url);
        }
        return await fetchBackend(url);
    } catch (ex) {
        console.log('Exception:', ex);
        return [500, errorHead(ex)];
    }
}

async function fetchBackend(url: string): Promise<[number, string]> {
    const resp = await axios.get(url);
    if (resp.status != 200) {
        console.log("ERROR: Non-OK status", resp.status, resp.statusText, url);
        return [resp.status, ""];
    }
    let text: string = "";
    try {
        text = JSON.stringify(resp.data);
    } catch (error) {
        text = resp.data;
    }
    return [resp.status, text]
}

async function fetchBrowser(url: string): Promise<[number, string]> {
    let resp = await fetch(url);
    if (resp.status != 200) {
        console.log("ERROR: Non-OK status", resp.status, resp.statusText, url);
        return [resp.status, ""];
    }
    const text: string = await resp.text();
    return [resp.status, text];
}

export async function httpPostFromBrowser(url: string, data: any): Promise<[number, any]> {
    const options = {
        method: 'POST',
        body: JSON.stringify(data),
        headers: {
            'Content-Type': 'application/json'
        }
    }
    //console.log('options.body', options.body);
    const resp = await fetch(url, options);
    if (resp.status != 200) {
        console.log("ERROR: Non-OK status", resp.status, resp.statusText, url);
        return [resp.status, ""];
    }
    const json: string = await resp.json();
    return [resp.status, json]
}
