import type { SessionMember } from "@/api/types/session.ts";
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
import { useUserStore } from "@/stores/user-store.tsx";
import { MemberDropdownMenu } from "./member-dropdown-menu";

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
