import { OTPInput } from "@/components/otp-input";
import { Button } from "@/components/ui/button.tsx";
import { Form, FormField } from "@/components/ui/form.tsx";
import { zodResolver } from "@hookform/resolvers/zod";
import { REGEXP_ONLY_DIGITS_AND_CHARS } from "input-otp";
import { useForm } from "react-hook-form";
import { z } from "zod";

const formSchema = z.object({
	otp: z.string().min(6).max(6).nonempty(),
});

export type OTPFormValues = z.infer<typeof formSchema>;

type UpdatePasswordOTPFormProps = {
	onCancel: () => void;
	onOtpCompleted: (otp: string) => void;
};

export function OTPForm({ onCancel, onOtpCompleted }: UpdatePasswordOTPFormProps) {
	const form = useForm<OTPFormValues>({
		resolver: zodResolver(formSchema),
		defaultValues: {
			otp: "",
		},
	});

	function handleSubmit(values: OTPFormValues) {
		onOtpCompleted(values.otp);
	}

	return (
		<Form {...form}>
			<form onSubmit={form.handleSubmit(handleSubmit)} className="space-y-4">
				<div className="flex justify-between items-end gap-2">
					<FormField
						name="otp"
						control={form.control}
						render={({ field }) => (
							<OTPInput
								length={6}
								pattern={REGEXP_ONLY_DIGITS_AND_CHARS}
								onComplete={(otp) => handleSubmit({ otp })}
								className="w-full"
								{...field}
							/>
						)}
					/>
					<Button type="submit" variant="secondary" disabled={form.formState.isSubmitting} onClick={onCancel}>
						Cancel
					</Button>
				</div>
			</form>
		</Form>
	);
}
