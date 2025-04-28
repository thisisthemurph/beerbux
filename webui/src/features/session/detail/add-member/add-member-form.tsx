import { Button } from "@/components/ui/button.tsx";
import { Form, FormControl, FormField, FormItem, FormLabel } from "@/components/ui/form.tsx";
import { Input } from "@/components/ui/input.tsx";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { z } from "zod";

const formSchema = z.object({
	username: z.string().nonempty(),
});

export type AddMemberFormValues = z.infer<typeof formSchema>;

type AddMemberFormProps = {
	onAdd: (username: string) => void;
};

export function AddMemberForm({ onAdd }: AddMemberFormProps) {
	const form = useForm<AddMemberFormValues>({
		resolver: zodResolver(formSchema),
		defaultValues: {
			username: "",
		},
	});

	function handleSubmit(values: AddMemberFormValues) {
		onAdd(values.username);
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
								<Input {...field} placeholder="Add a member by username" />
							</FormControl>
						</FormItem>
					)}
				/>

				<Button type="submit">Add</Button>
			</form>
		</Form>
	);
}
