import { OTPInput } from "@/components/otp-input";
import { Button } from "@/components/ui/button.tsx";
import { Form, FormControl, FormField, FormItem, FormLabel } from "@/components/ui/form.tsx";
import { Input } from "@/components/ui/input.tsx";
import { zodResolver } from "@hookform/resolvers/zod";
import { REGEXP_ONLY_DIGITS_AND_CHARS } from "input-otp";
import { useForm } from "react-hook-form";
import { z } from "zod";

const formSchema = z.object({
	password: z.string().nonempty("Password is required"),
});

const otpFormSchema = z.object({
	otp: z.string().min(6).max(6).nonempty(),
});

export type PasswordFormValues = z.infer<typeof formSchema>;
export type OTPFormValues = z.infer<typeof otpFormSchema>;

type UpdatePasswordFormProps = {
	onSuccess: (values: PasswordFormValues) => void;
};

export function UpdatePasswordForm({ onSuccess }: UpdatePasswordFormProps) {
	const form = useForm<PasswordFormValues>({
		resolver: zodResolver(formSchema),
		defaultValues: {
			password: "",
		},
	});

	function handleSubmit(values: PasswordFormValues) {
		console.log(values);
		onSuccess(values);
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

type UpdatePasswordOTPFormProps = {
	onCancel: () => void;
	onSuccess: (otp: string) => void;
};

export function UpdatePasswordOTPForm({ onCancel, onSuccess }: UpdatePasswordOTPFormProps) {
	const form = useForm<OTPFormValues>({
		resolver: zodResolver(otpFormSchema),
		defaultValues: {
			otp: "",
		},
	});

	function handleSubmit(values: OTPFormValues) {
		console.log("submitted otp", values);
		// TODO: Verify OTP and reset password
		onSuccess(values.otp);
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
