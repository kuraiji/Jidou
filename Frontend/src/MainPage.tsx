//import {Message} from "@/app/actions";
import MyForm from "./MyForm";

export interface Message {
    name: string;
    message: string;
    date: string;
}

export default function MainPage({message}: {message: Message[]}) {
    return (
        <div className="flex flex-col items-center justify-center w-full max-w-screen">
            <h1 className="text-4xl mb-2">Jidou</h1>
            <h2 className="text-xl mb-2">Showing the most recent 10 messages.</h2>
            <ul className="flex flex-col items-center justify-center w-full max-w-screen gap-1 mb-5">
                <DisplayMessage message={{name: "Name", message: "Message", date: "Date"}} isHeader={true}/>
                {
                    message.map((item, i) => (
                        <DisplayMessage message={item} key={i} />
                    ))
                }
            </ul>
            <MyForm/>
        </div>
    );
}

function DisplayMessage(props: { message: Message, isHeader?: boolean }) {
    return (
        <li className={`grid grid-cols-6 w-full max-w-screen gap-2 border-b`}>
            <div className="flex flex-col items-center justify-center border-r">
                <p className={`${props.isHeader ? "font-bold" : ""}`}>{props.message.name}</p>
            </div>
            <div className={`col-span-4 border-r ${props.isHeader ? "flex flex-col items-center justify-center" : ""}`}>
                <p className={`${props.isHeader ? "font-bold" : ""}`}>{props.message.message}</p>
            </div>
            <div className="col-start-6 flex flex-col items-center justify-center">
                {
                    props.isHeader ?
                        <p className={"font-bold"} >{props.message.date}</p>
                        :
                        <p>{new Date(props.message.date).toLocaleString()}</p>
                }
            </div>
        </li>
    )
}