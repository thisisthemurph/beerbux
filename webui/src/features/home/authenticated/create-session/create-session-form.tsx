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
	name: z
		.string()
		.nonempty()
		.min(2, "The name must be at least 2 characters")
		.max(25, "The name cannot be more than 25 characters"),
});

export type CreateSessionFormValues = z.infer<typeof formSchema>;

type NewSessionFormProps = {
	onCreate: (values: CreateSessionFormValues) => void;
};

export function CreateSessionForm({ onCreate }: NewSessionFormProps) {
	const form = useForm<CreateSessionFormValues>({
		resolver: zodResolver(formSchema),
		defaultValues: {
			name: "",
		},
	});

	function handleSubmit(values: CreateSessionFormValues) {
		onCreate(values);
	}

	return (
		<Form {...form}>
			<form onSubmit={form.handleSubmit(handleSubmit)} className="space-y-4">
				<FormField
					name="name"
					control={form.control}
					render={({ field }) => (
						<FormItem>
							<FormLabel htmlFor={field.name}>Session name</FormLabel>
							<FormControl>
								<Input
									type="text"
									{...field}
									placeholder="Name of your session"
								/>
							</FormControl>
						</FormItem>
					)}
				/>

				<Button type="submit">Create</Button>
			</form>
		</Form>
	);
}
