import { Button } from "@/components/ui/button.tsx";
import { Form, FormControl, FormField, FormItem, FormLabel } from "@/components/ui/form.tsx";
import { Input } from "@/components/ui/input.tsx";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { z } from "zod";

const formSchema = z.object({
	email: z.string().email("Invalid email address").nonempty("Email is required"),
});

export type UpdateEmailFormValues = z.infer<typeof formSchema>;

type UpdateEmailFormProps = {
	onSubmit: (email: string) => void;
};

export function UpdateEmailForm({ onSubmit }: UpdateEmailFormProps) {
	const form = useForm<UpdateEmailFormValues>({
		resolver: zodResolver(formSchema),
		defaultValues: {
			email: "",
		},
	});

	function handleSubmit(values: UpdateEmailFormValues) {
		onSubmit(values.email);
	}

	return (
		<Form {...form}>
			<form onSubmit={form.handleSubmit(handleSubmit)} className="space-y-4">
				<div className="flex justify-between items-end gap-2">
					<FormField
						name="email"
						control={form.control}
						render={({ field }) => (
							<FormItem className="w-full">
								<FormLabel htmlFor={field.name}>Email</FormLabel>
								<FormControl>
									<Input type="email" {...field} />
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
