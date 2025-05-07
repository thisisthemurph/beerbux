import type { SessionMember } from "@/api/types/session.ts";
import { InformationButton } from "@/components/information-button.tsx";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card.tsx";
import { UserAvatar } from "@/components/user-avatar.tsx";
import { useInformationDialog } from "@/hooks/use-information-dialog.tsx";
import { useUserAvatarDataBySession } from "@/hooks/user-avatar-data.ts";
import { cn } from "@/lib/utils";
import { useUserStore } from "@/stores/user-store.tsx";
import { MemberDropdownMenu } from "./member-dropdown-menu";

type MemberDetailsCardProps = {
	sessionId: string;
	members: SessionMember[];
	showMemberDropdownMenu: boolean;
	onChangeMemberAdminState: (memberId: string, newAdminState: boolean) => void;
	onRemoveMember: (memberId: string) => void;
	onSelectMember: (memberId: string) => void;
};

export function MemberDetailsCard({
	sessionId,
	members,
	showMemberDropdownMenu,
	onChangeMemberAdminState,
	onRemoveMember,
	onSelectMember,
}: MemberDetailsCardProps) {
	const [openInformationDialog, InformationDialog] = useInformationDialog();
	const user = useUserStore((state) => state.user);
	const avatarData = useUserAvatarDataBySession(sessionId);

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
						<section key={m.id} className="grid grid-cols-5 hover:bg-muted">
							<button
								type="button"
								className="col-span-4 flex w-full items-center gap-6 pl-6 cursor-pointer"
								onClick={() => onSelectMember(m.id)}
							>
								<UserAvatar data={avatarData[m.username]} variant={m.isAdmin ? "prominent" : "default"} />
								<div className="w-full py-6 text-left tracking-wider font-semibold">
									<span>{m.username}</span>
								</div>
							</button>
							<div className="flex items-center gap-2">
								<Balance {...m.transactionSummary} />
								{m.id !== user?.id && showMemberDropdownMenu && (
									<MemberDropdownMenu
										username={m.username}
										isAdmin={m.isAdmin}
										onChangeAdminState={() => onChangeMemberAdminState(m.id, !m.isAdmin)}
										onRemoveMember={() => onRemoveMember(m.id)}
									/>
								)}
							</div>
						</section>
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
