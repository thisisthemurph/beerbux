import type { SessionMember } from "@/api/types/session.ts";
import {
	Drawer,
	DrawerContent,
	DrawerDescription,
	DrawerHeader,
	DrawerTitle,
} from "@/components/ui/drawer.tsx";
import { CreateTransactionForm } from "@/features/session/detail/create-transaction/careate-transaction-form.tsx";
import type { DrawerToggleProps } from "@/types.ts";
import { Beer } from "lucide-react";
import { useState } from "react";

interface CreateTransactionDrawerProps extends DrawerToggleProps {
	members: SessionMember[];
	onTransactionCreate: (values: Record<string, number>) => void;
}

export function CreateTransactionDrawer({
	members,
	onTransactionCreate,
	...drawerToggleProps
}: CreateTransactionDrawerProps) {
	const [total, setTotal] = useState(members.length);

	return (
		<Drawer {...drawerToggleProps}>
			<DrawerContent>
				<DrawerHeader>
					<section className="flex justify-between items-start">
						<DrawerTitle>Buy a round</DrawerTitle>
						<p className="flex items-center gap-1 font-semibold text-lg">
							{total} <Beer className="size-4" />
						</p>
					</section>
					<DrawerDescription>Select the members you want to buy a round for.</DrawerDescription>
				</DrawerHeader>
				<section className="p-4">
					<CreateTransactionForm
						members={members}
						onTotalChanged={(amount) => setTotal(amount)}
						onTransactionCreate={onTransactionCreate}
					/>
				</section>
			</DrawerContent>
		</Drawer>
	);
}
