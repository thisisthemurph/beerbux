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

const formSchema = z
	.object({
		name: z
			.string()
			.nonempty()
			.max(50, "Your name cannot be more than 50 characters"),
		username: z
			.string()
			.nonempty()
			.min(3, "Your username must be at least 3 characters")
			.max(25, "Your username cannot be more than 25 characters"),
		password: z
			.string()
			.nonempty()
			.min(8, "Your password must be at least 8 characters"),
		verificationPassword: z.string().nonempty(),
	})
	.refine((schema) => schema.password === schema.verificationPassword, {
		message: "Passwords do not match",
		path: ["verificationPassword"],
	});

export type SignupFormValues = z.infer<typeof formSchema>;

interface LoginFormProps {
	onSubmit: (values: SignupFormValues) => void;
}

export function SignupForm({ onSubmit }: LoginFormProps) {
	const form = useForm<SignupFormValues>({
		resolver: zodResolver(formSchema),
		defaultValues: {
			name: "",
			username: "",
			password: "",
			verificationPassword: "",
		},
	});

	function handleSubmit(values: SignupFormValues) {
		onSubmit(values);
	}

	return (
		<Form {...form}>
			<form onSubmit={form.handleSubmit(handleSubmit)} className="space-y-4">
				<FormField
					name="name"
					control={form.control}
					render={({ field }) => (
						<FormItem>
							<FormLabel htmlFor={field.name}>Name</FormLabel>
							<FormControl>
								<Input {...field} />
							</FormControl>
						</FormItem>
					)}
				/>
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

				<FormField
					name="verificationPassword"
					control={form.control}
					render={({ field }) => (
						<FormItem>
							<FormLabel htmlFor={field.name}>Confirm password</FormLabel>
							<FormControl>
								<Input type="password" {...field} />
							</FormControl>
						</FormItem>
					)}
				/>

				<Button type="submit">Sign up</Button>
			</form>
		</Form>
	);
}
