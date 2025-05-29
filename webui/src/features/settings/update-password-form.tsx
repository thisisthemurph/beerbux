import { Button } from "@/components/ui/button.tsx";
import { Form, FormControl, FormField, FormItem, FormLabel } from "@/components/ui/form.tsx";
import { Input } from "@/components/ui/input.tsx";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { z } from "zod";

const formSchema = z.object({
	password: z.string().nonempty("Password is required"),
});

export type PasswordFormValues = z.infer<typeof formSchema>;

type UpdatePasswordFormProps = {
	onSubmit: (password: string) => void;
};

export function UpdatePasswordForm({ onSubmit }: UpdatePasswordFormProps) {
	const form = useForm<PasswordFormValues>({
		resolver: zodResolver(formSchema),
		defaultValues: {
			password: "",
		},
	});

	function handleSubmit(values: PasswordFormValues) {
		onSubmit(values.password);
	}

	return (
		<Form {...form}>
			<form onSubmit={form.handleSubmit(handleSubmit)} className="space-y-4">
				<div className="flex justify-between items-end gap-2">
					<FormField
						name="password"
						control={form.control}
						render={({ field }) => (
							<FormItem className="w-full">
								<FormLabel htmlFor={field.name}>Password</FormLabel>
								<FormControl>
									<Input type="password" {...field} />
								</FormControl>
							</FormItem>
						)}
					/>
					<Button
						type="submit"
						variant="secondary"
						disabled={!form.formState.isValid || form.formState.isSubmitting}
					>
						Request
					</Button>
				</div>
			</form>
		</Form>
	);
}
