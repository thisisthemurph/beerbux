import type { SessionMember } from "@/api/types";
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
import { Minus, Plus } from "lucide-react";
import { useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { toast } from "sonner";
import { z } from "zod";

const formSchema = z.record(
	z.string(),
	z.number().int().min(0, "Amount cannot be less than 0"),
);

type TransactionFormValues = z.infer<typeof formSchema>;
export type Transaction = Record<string, number>;

type TransactionFormProps = {
	members: SessionMember[];
	onTransactionCreate: (values: Record<string, number>) => void;
	onTotalChanged: (amount: number) => void;
};

export function TransactionForm({
	members,
	onTransactionCreate,
	onTotalChanged,
}: TransactionFormProps) {
	const [total, setTotal] = useState(members.length);
	const form = useForm<TransactionFormValues>({
		resolver: zodResolver(formSchema),
		defaultValues: Object.fromEntries(members.map((m) => [m.id, 1])),
	});

	useEffect(() => {
		onTotalChanged(total);
	}, [total, onTotalChanged]);

	function handleSubmit(values: TransactionFormValues) {
		const transaction = Object.entries(values).reduce<Transaction>(
			(acc, [memberId, amount]) => {
				const value = Number(amount);
				if (value > 0) acc[memberId] = value;
				return acc;
			},
			{},
		);

		if (total <= 0) {
			toast.error("The transaction must have at least one non-zero value.");
			return;
		}

		onTransactionCreate(transaction);
	}

	function updateTotal() {
		setTotal(Object.values(form.getValues()).reduce((sum, v) => sum + v, 0));
	}

	return (
		<div>
			<Form {...form}>
				<form onSubmit={form.handleSubmit(handleSubmit)} className="space-y-4">
					{members.map((member) => (
						<FormField
							key={member.id}
							name={member.id}
							control={form.control}
							render={({ field }) => (
								<FormItem className="flex justify-between items-center">
									<FormLabel htmlFor={field.name} className="w-full">
										{member.name}
									</FormLabel>
									<Button
										type="button"
										size="icon"
										variant="secondary"
										className="rounded-full"
										onClick={(e) => {
											e.preventDefault();
											form.setValue(member.id, updateValueBy(field.value, -1));
											updateTotal();
										}}
									>
										<Minus />
									</Button>
									<FormControl className="w-1/3">
										<Input
											{...field}
											inputMode="numeric"
											className="text-center font-semibold"
											onChange={(e) => {
												const value = Number(e.target.value);
												form.setValue(
													member.id,
													Number.isNaN(value) ? 0 : value,
													{ shouldValidate: true },
												);

												updateTotal();
											}}
										/>
									</FormControl>
									<Button
										type="button"
										size="icon"
										variant="secondary"
										className="rounded-full"
										onClick={(e) => {
											e.preventDefault();
											form.setValue(member.id, updateValueBy(field.value, 1));
											updateTotal();
										}}
									>
										<Plus />
									</Button>
								</FormItem>
							)}
						/>
					))}
					<Button type="submit">Create transaction</Button>
				</form>
			</Form>
		</div>
	);
}

function updateValueBy(value: string | number, amount: number): number {
	const newValue =
		(typeof value === "string" ? Number.parseInt(value) : value) + amount;
	if (newValue < 0) {
		return 0;
	}
	return newValue;
}
