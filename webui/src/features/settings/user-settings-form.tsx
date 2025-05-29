import { Button } from "@/components/ui/button.tsx";
import { Form, FormControl, FormField, FormItem, FormLabel } from "@/components/ui/form.tsx";
import { Input } from "@/components/ui/input.tsx";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { z } from "zod";

const formSchema = z.object({
	name: z.string().nonempty(),
	username: z.string().nonempty(),
});

export type UserSettingsFormValues = z.infer<typeof formSchema>;

type UserSettingsFormProps = {
	name: string;
	username: string;
	onSubmit: (values: UserSettingsFormValues) => void;
};

export function UserSettingsForm({ name, username, onSubmit }: UserSettingsFormProps) {
	const form = useForm<UserSettingsFormValues>({
		resolver: zodResolver(formSchema),
		defaultValues: {
			name: name || "",
			username: username || "",
		},
	});

	function handleSubmit(values: UserSettingsFormValues) {
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

				<Button
					type="submit"
					variant="secondary"
					className="w-full"
					disabled={!form.formState.isValid || !form.formState.isDirty || form.formState.isSubmitting}
				>
					Update
				</Button>
			</form>
		</Form>
	);
}
