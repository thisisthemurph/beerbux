import { PageHeading } from "@/components/page-heading.tsx";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { UpdateEmailForm } from "@/features/settings/update-email-form.tsx";
import { UpdatePasswordForm, UpdatePasswordOTPForm } from "@/features/settings/update-password-form.tsx";
import { UserSettingsForm } from "@/features/settings/user-settings-form.tsx";
import { useState } from "react";

export default function SettingsPage() {
	const [showPasswordOTP, setShowPasswordOTP] = useState(false);

	return (
		<>
			<PageHeading title="Settings" />
			<section className="space-y-4">
				<Card>
					<CardHeader>
						<CardTitle>User settings</CardTitle>
						<CardDescription className="sr-only">user settings</CardDescription>
					</CardHeader>
					<CardContent>
						<UserSettingsForm />
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
