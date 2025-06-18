'use server'
import { SSMClient, GetParameterCommand } from "@aws-sdk/client-ssm";
import axios from "axios";

export interface Message {
    name: string;
    message: string;
    date: string;
}

interface ApiReq {
    url: string;
    authorization: string;
}

export async function fetchApiUrl(): Promise<ApiReq> {
    try {
        const client = new SSMClient({region: process.env.REGION});
        const GetIpCommand = new GetParameterCommand({
            Name: "/JIDOU-API/BACKEND_IP"
        });
        const GetPortCommand = new GetParameterCommand({
            Name: "/JIDOU-API/BACKEND_PORT"
        });
        const GetAuthCommand = new GetParameterCommand({
            Name: "/JIDOU-API/EC2_KEY"
        })
        const ip_res = await client.send(GetIpCommand);
        const port_res = await client.send(GetPortCommand);
        const auth_res = await client.send(GetAuthCommand);
        if(!ip_res.Parameter?.Value || !port_res.Parameter?.Value || !auth_res.Parameter?.Value)
            throw new Error("Couldn't Get Parameter Values");
        return {
            url: `${ip_res.Parameter.Value}:${port_res.Parameter.Value}`,
            authorization: auth_res.Parameter.Value
        }
    }
    catch (error) {
        throw new Error(`Failed to fetch api url - ${error}`);
    }
}

export async function fetchMessages(): Promise<Message[]> {
    try {
        const apiReq = await fetchApiUrl();
        return await axios.get(`http://${apiReq.url}/`,
            {
                headers: {
                    "Authorization": apiReq.authorization
                }
            }).then(response => {
                if (response.status >= 400 && response.status <= 499) throw new Error(response.statusText);
                return response.data;
        });
    }
    catch(error) {
        throw new Error(`Failed to fetch messages - ${error}`);
    }
}

export async function postMyMessage(name: string, message: string): Promise<void> {
    try {
        const apiReq = await fetchApiUrl();
        return await axios.post(`http://${apiReq.url}/`,
            {
                message: message,
                name: name,
            },
            {
                headers: {
                    "Authorization": apiReq.authorization,
                    "Content-Type": "application/json",
                }
            }).then(response => {
            if (response.status >= 400 && response.status <= 499) throw new Error(response.statusText);
        });
    }
    catch(error) {
        throw new Error(`Failed to post message - could be due to server error or your name/message has profanity in it.`);
    }
}