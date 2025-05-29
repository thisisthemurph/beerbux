import { InputOTP, InputOTPGroup, InputOTPSlot } from "@/components/ui/input-otp";
import { cn } from "@/lib/utils";
import type * as React from "react";
import { useState } from "react";

type OTPInputWrapperProps = Omit<
	React.ComponentProps<typeof InputOTP>,
	"children" | "render" | "maxLength" | "value" | "onChange"
> & {
	length?: number;
	onComplete?: (value: string) => void;
};

export function OTPInput({
	length = 6,
	onComplete,
	className,
	containerClassName,
	...rest
}: OTPInputWrapperProps) {
	const [value, setValue] = useState("");

	return (
		<InputOTP
			{...rest}
			maxLength={length}
			value={value}
			onChange={(v) => {
				setValue(v);
				if (onComplete && v.length === length) {
					onComplete(v);
				}
			}}
			className={cn("w-full", className)}
			containerClassName={cn("", containerClassName)}
		>
			<InputOTPGroup>
				{Array.from({ length }).map((_, index) => (
					// biome-ignore lint/suspicious/noArrayIndexKey: <explanation>
					<InputOTPSlot key={index} index={index} />
				))}
			</InputOTPGroup>
		</InputOTP>
	);
}
