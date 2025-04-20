import type { SessionMember } from "@/api/types.ts";
import { InformationButton } from "@/components/information-button.tsx";
import {
	Card,
	CardContent,
	CardHeader,
	CardTitle,
} from "@/components/ui/card.tsx";
import { UserAvatar } from "@/components/user-avatar.tsx";
import { useInformationDialog } from "@/hooks/use-information-dialog.tsx";
import type { AvatarData } from "@/hooks/user-avatar-data.ts";
import { cn } from "@/lib/utils";
import {
	DropdownMenu,
	DropdownMenuContent,
	DropdownMenuItem,
	DropdownMenuLabel,
	DropdownMenuSeparator,
	DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {
	Crown,
	EllipsisVertical,
	ShieldMinus,
	ShieldUser,
	UserMinus,
} from "lucide-react";
import { Button } from "@/components/ui/button.tsx";
import { useUserStore } from "@/stores/user-store.tsx";

type MemberDetailsCardProps = {
	members: SessionMember[];
	avatarData: Record<string, AvatarData>;
	showMemberDropdownMenu: boolean;
	onChangeMemberAdminState: (memberId: string, newAdminState: boolean) => void;
	onRemoveMember: (memberId: string) => void;
};

export function MemberDetailsCard({
	members,
	avatarData,
	showMemberDropdownMenu,
	onChangeMemberAdminState,
	onRemoveMember,
}: MemberDetailsCardProps) {
	const [openInformationDialog, InformationDialog] = useInformationDialog();
	const user = useUserStore((state) => state.user);

	const handleInformationClick = () => {
		openInformationDialog({
			title: "Members",
			description:
				"This section shows a list of all members in the session and a summary of how many beers that person owes or is owed.",
		});
	};

	return (
		<>
			<InformationDialog />
			<Card>
				<CardHeader>
					<section className="flex items-center justify-between">
						<CardTitle>Members</CardTitle>
						<InformationButton onClick={handleInformationClick} />
					</section>
				</CardHeader>
				<CardContent className="px-0">
					{members.map((m) => (
						<div
							key={m.id}
							className="group flex items-center gap-4 px-6 hover:bg-muted"
						>
							<UserAvatar
								data={avatarData[m.username]}
								tooltip={m.name}
								variant={m.isAdmin ? "prominent" : "default"}
							/>
							<div className="flex justify-between items-center w-full">
								<button
									type="button"
									className="flex items-center gap-2 py-6 font-semibold"
								>
									<span>{m.username}</span>
								</button>
								<div
									className={cn(
										"",
										showMemberDropdownMenu && "w-16 text-left grid grid-cols-2",
									)}
								>
									<Balance {...m.transactionSummary} />
									{m.id !== user?.id && showMemberDropdownMenu && (
										<MemberDropdownMenu
											username={m.username}
											isAdmin={m.isAdmin}
											onChangeAdminState={() =>
												onChangeMemberAdminState(m.id, !m.isAdmin)
											}
											onRemoveMember={() => onRemoveMember(m.id)}
										/>
									)}
								</div>
							</div>
						</div>
					))}
				</CardContent>
			</Card>
		</>
	);
}

type MemberDropdownMenuProps = {
	username: string;
	isAdmin: boolean;
	onChangeAdminState: () => void;
	onRemoveMember: () => void;
};

function MemberDropdownMenu({
	username,
	isAdmin,
	onChangeAdminState,
	onRemoveMember,
}: MemberDropdownMenuProps) {
	return (
		<DropdownMenu>
			<DropdownMenuTrigger
				className="py-6 font-semibold w-full text-left"
				asChild
			>
				<Button variant="ghost" size="icon" className="rounded-full">
					<EllipsisVertical className="size-6 text-muted-foreground group-hover:text-primary" />
				</Button>
			</DropdownMenuTrigger>
			<DropdownMenuContent className="min-w-46 mx-4">
				<DropdownMenuLabel className="flex items-center justify-between font-semibold text-lg">
					<span>{username}</span>
					{isAdmin && (
						<span>
							<Crown className="size-4" />
						</span>
					)}
				</DropdownMenuLabel>
				<DropdownMenuSeparator />
				<DropdownMenuItem
					className="text-lg gap-4 cursor-pointer"
					onClick={onChangeAdminState}
				>
					{isAdmin ? (
						<ShieldMinus className="size-6" />
					) : (
						<ShieldUser className="size-6" />
					)}
					<span>{isAdmin ? "Remove admin" : "Promote to admin"}</span>
				</DropdownMenuItem>
				<DropdownMenuItem
					className="text-lg gap-4 cursor-pointer"
					onClick={onRemoveMember}
				>
					<UserMinus className="size-6" />
					<span>Remove from session</span>
				</DropdownMenuItem>
			</DropdownMenuContent>
		</DropdownMenu>
	);
}

function Balance({ credit, debit }: { credit: number; debit: number }) {
	const value = credit - debit;

	return (
		<p
			className={cn(
				"font-semibold text-muted-foreground text-lg flex items-center",
				value > 0 && "text-green-600",
				value < 0 && "text-red-600",
			)}
		>
			${value < 0 ? value * -1 : value}
		</p>
	);
}
