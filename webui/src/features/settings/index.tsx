import useUserClient from "@/api/user-client.ts";
import { PageHeading } from "@/components/page-heading.tsx";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { UpdateEmailForm } from "@/features/settings/update-email-form.tsx";
import { UpdatePasswordForm, UpdatePasswordOTPForm } from "@/features/settings/update-password-form.tsx";
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
	const { updateUser } = useUserClient();

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
						<UpdateEmailForm />
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
							<UpdatePasswordOTPForm
								onCancel={() => setShowPasswordOTP(false)}
								onSuccess={() => setShowPasswordOTP(false)}
							/>
						) : (
							<UpdatePasswordForm onSuccess={() => setShowPasswordOTP(true)} />
						)}
					</CardContent>
				</Card>
			</section>
		</>
	);
}
