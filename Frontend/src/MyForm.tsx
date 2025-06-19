'use client'
import { Field } from '@base-ui-components/react/field';
import { Form } from '@base-ui-components/react/form';
import styles from './index.module.css'
import { z } from 'zod';
import React from "react";
import {useMutation} from "@tanstack/react-query";
import axios from "axios";
import {api_url} from "./App.tsx";

const schema = z.object({
    name: z.string().min(2).max(20),
    message: z.string().min(2).max(60),
})

export default function MyForm() {
    const [errors, setErrors] = React.useState({});
    const [disabled, setDisabled] = React.useState(false);

    const mutation = useMutation({
        mutationFn: async (newMessage: {name: string; message: string}) => {
            try {
                return await axios.post(`${api_url}`,
                    {
                        message: newMessage.message,
                        name: newMessage.name,
                    },
                    {
                        headers: {
                            "Content-Type": "application/json",
                        }
                    }).then(response => {
                    if (response.status >= 400 && response.status <= 499) throw new Error(response.statusText);
                });
            }
            catch {
                throw new Error(`Failed to post message - could be due to server error or your name/message has profanity in it.`);
            }
        }
    })

    async function submitForm(event: React.FormEvent<HTMLFormElement>) {
        event.preventDefault();
        setDisabled(true);

        const formData = new FormData(event.currentTarget);
        const result = schema.safeParse(Object.fromEntries(formData as any));

        if (!result.success) {
            setDisabled(false);
            return {
                errors: result.error.flatten().fieldErrors,
            };
        }
        try {
            await mutation.mutateAsync({name: result.data.name, message: result.data.message})
        }
        catch (error) {
            setDisabled(false);
            return {
                errors: {message: `${error}`},
            }
        }
        setDisabled(false);
        const e = document.getElementById('message');
        if (e instanceof HTMLInputElement) e.value = '';
        location.reload();
        return {
            errors: {}
        };
    }

    return (
        <Form
            className={styles.Form}
            errors={errors}
            onClearErrors={setErrors}
            onSubmit={async (event) => {
                const response = await submitForm(event);
                setErrors(response.errors);
            }}
        >
            <Field.Root name="name" className={styles.Field}>
                <Field.Label className={styles.Label}>Name</Field.Label>
                <Field.Control id="name" required placeholder="Enter name" className={styles.Input} />
                <Field.Error className={styles.Error} />
            </Field.Root>
            <Field.Root name="message" className={styles.Field}>
                <Field.Label className={styles.Label}>Message</Field.Label>
                <Field.Control id="message" required type="text" placeholder="Enter message" className={styles.Input} />
                <Field.Error className={styles.Error} />
            </Field.Root>
            <button type="submit" className={styles.Button} disabled={disabled}>
                Submit
            </button>
        </Form>
    );
}