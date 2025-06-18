'use client'
import { Field } from '@base-ui-components/react/field';
import { Form } from '@base-ui-components/react/form';
import styles from '@/app/index.module.css'
import { z } from 'zod';
import React from "react";
import {postMyMessage} from "@/app/actions";
import { useRouter } from 'next/navigation';

const schema = z.object({
    name: z.string().min(2).max(20),
    message: z.string().min(2).max(60),
})

export default function MyForm() {
    const router = useRouter();
    const [errors, setErrors] = React.useState({});
    const [disabled, setDisabled] = React.useState(false);

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
            await postMyMessage(result.data.name, result.data.message);
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
        router.push('/');
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