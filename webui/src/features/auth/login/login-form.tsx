import { Button } from "@/components/ui/button.tsx";
import {
	Form,
	FormControl,
	FormField,
	FormItem,
	FormLabel,
} from "@/components/ui/form.tsx";
import { Input } from "@/components/ui/input.tsx";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { z } from "zod";

const formSchema = z.object({
	username: z.string().nonempty(),
	password: z.string().nonempty(),
});

export type LoginFormValues = z.infer<typeof formSchema>;

interface LoginFormProps {
	onSubmit: (values: LoginFormValues) => void;
}

export function LoginForm({ onSubmit }: LoginFormProps) {
	const form = useForm<LoginFormValues>({
		resolver: zodResolver(formSchema),
		defaultValues: {
			username: "",
			password: "",
		},
	});

	function handleSubmit(values: LoginFormValues) {
		onSubmit(values);
	}

	return (
		<Form {...form}>
			<form onSubmit={form.handleSubmit(handleSubmit)} className="space-y-4">
				<FormField
					name="username"
					control={form.control}
					render={({ field }) => (
						<FormItem>
							<FormLabel htmlFor={field.name}>Username</FormLabel>
							<FormControl>
								<Input {...field} />
							</FormControl>
						</FormItem>
					)}
				/>

				<FormField
					name="password"
					control={form.control}
					render={({ field }) => (
						<FormItem>
							<FormLabel htmlFor={field.name}>Password</FormLabel>
							<FormControl>
								<Input type="password" {...field} />
							</FormControl>
						</FormItem>
					)}
				/>

				<Button type="submit">Login</Button>
			</form>
		</Form>
	);
}
