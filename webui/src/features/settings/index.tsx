import useAuthClient from "@/api/auth-client.ts";
import useUserClient from "@/api/user-client.ts";
import { PageHeading } from "@/components/page-heading.tsx";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { OTPForm } from "@/features/settings/otp-form.tsx";
import { UpdateEmailForm } from "@/features/settings/update-email-form.tsx";
import { UpdatePasswordForm } from "@/features/settings/update-password-form.tsx";
import { UserSettingsForm } from "@/features/settings/user-settings-form.tsx";
import type { UserSettingsFormValues } from "@/features/settings/user-settings-form.tsx";
import { tryCatch } from "@/lib/try-catch.ts";
import { useUserStore } from "@/stores/user-store.tsx";
import { useState } from "react";
import { toast } from "sonner";

export default function SettingsPage() {
	const user = useUserStore((state) => state.user);
	const setUser = useUserStore((state) => state.setUser);
	const [showPasswordOTP, setShowPasswordOTP] = useState(false);
	const [showEmailOTP, setShowEmailOTP] = useState(false);
	const { updateUser } = useUserClient();
	const { initializePasswordReset, resetPassword, initializeEmailUpdate, updateEmail } = useAuthClient();

	async function handleUpdateUser(values: UserSettingsFormValues) {
		if (!user) return;
		const { data, err } = await tryCatch(updateUser(values.username, values.name));
		if (err) {
			toast.error("Failed to update user details", {
				description: err.message,
			});
			return;
		}

		setUser({ ...user, ...data });
		toast.success("User details updated successfully", {
			description:
				data.username !== user.username
					? "Next time you log in, you will need to log in with your new username."
					: undefined,
		});
	}

	async function handleInitializePasswordUpdate(password: string) {
		const { err } = await tryCatch(initializePasswordReset(password));
		if (err) {
			toast.error("Failed to initialize password reset", {
				description: err.message,
			});
			return;
		}

		setShowPasswordOTP(true);
		toast.success("Password reset requested. Please check your email for the OTP.");
	}

	async function handleUpdatePassword(otp: string) {
		const { err } = await tryCatch(resetPassword(otp));
		if (err) {
			toast.error("Failed to reset password", {
				description: err.message,
			});
			return;
		}
		setShowPasswordOTP(false);
		toast.success("Password reset successfully. You can now log in with your new password.");
	}

	async function handleInitializeEmailUpdate(newEmail: string) {
		const { err } = await tryCatch(initializeEmailUpdate(newEmail));
		if (err) {
			toast.error("Failed to start the email update process", {
				description: err.message,
			});
			return;
		}

		setShowEmailOTP(true);
		toast.success("Email update requested. Please check your email for the OTP.");
	}

	async function handleUpdateEmail(otp: string) {
		const { err } = await tryCatch(updateEmail(otp));
		if (err) {
			toast.error("Could not update your email address", {
				description: err.message,
			});
			return;
		}
		setShowEmailOTP(false);
		toast.success("Email updated successfully.");
	}

	if (!user) {
		return <div className="p-4">You must be logged in to view this page.</div>;
	}

	return (
		<>
			<PageHeading title="Settings" />
			<section className="space-y-4">
				<Card>
					<CardHeader>
						<CardTitle>User details</CardTitle>
						<CardDescription className="sr-only">user settings</CardDescription>
					</CardHeader>
					<CardContent>
						<UserSettingsForm onSubmit={handleUpdateUser} name={user.name} username={user.username} />
					</CardContent>
				</Card>
				<Card>
					<CardHeader>
						<CardTitle>Email address</CardTitle>
						<CardDescription>
							Change your email address. You will receive an email with a verification code.
						</CardDescription>
					</CardHeader>
					<CardContent>
						{showEmailOTP ? (
							<OTPForm
								onCancel={() => setShowEmailOTP(false)}
								onOtpCompleted={(otp) => handleUpdateEmail(otp)}
							/>
						) : (
							<UpdateEmailForm onSubmit={(newEmail) => handleInitializeEmailUpdate(newEmail)} />
						)}
					</CardContent>
				</Card>
				<Card>
					<CardHeader>
						<CardTitle>Reset password</CardTitle>
						<CardDescription>
							{showPasswordOTP
								? "Enter the OTP sent to your email to reset your password."
								: "Enter your new password and request an OTP to reset your password."}
						</CardDescription>
					</CardHeader>
					<CardContent>
						{showPasswordOTP ? (
							<OTPForm
								onCancel={() => setShowPasswordOTP(false)}
								onOtpCompleted={(otp) => handleUpdatePassword(otp)}
							/>
						) : (
							<UpdatePasswordForm onSubmit={(newPassword) => handleInitializePasswordUpdate(newPassword)} />
						)}
					</CardContent>
				</Card>
			</section>
		</>
	);
}
