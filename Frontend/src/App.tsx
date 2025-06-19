import './App.css'
import {useQuery} from "@tanstack/react-query";
import axios from "axios";
import MainPage, {type Message} from "./MainPage.tsx";
import {useState} from "react";

export const api_url = import.meta.env.MODE === "development" ? "http://127.0.0.1:1323/api" : "/api";

function App() {
    const [errorMessage, setErrorMessage] = useState<String | null>(null);
    const {data: messages, isLoading} = useQuery<Message[]>({
        queryKey: ["messages"],
        queryFn: async () => {
            try {
                return await axios.get(`${api_url}`,
                    {
                    }).then(response => {
                    if (response.status >= 400 && response.status <= 499) throw new Error(response.statusText);
                    return response.data;
                });
            }
            catch(error) {
                setErrorMessage(`Failed to fetch messages - ${error}`);
            }
        }
    });
    return (
        <>
            {
                isLoading || !messages || messages.length === 0 ?
                    <>
                        <h1>Loading...</h1>
                        {
                            errorMessage ?
                                <p>{errorMessage}</p>
                            : null
                        }
                    </>
                :
                    <MainPage message={messages}/>
            }
        </>
    )
}

export default App
