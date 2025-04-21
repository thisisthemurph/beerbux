import {
	DropdownMenu,
	DropdownMenuContent,
	DropdownMenuItem,
	DropdownMenuLabel,
	DropdownMenuSeparator,
	DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu.tsx";
import { Button } from "@/components/ui/button.tsx";
import {
	Crown,
	EllipsisVertical,
	ShieldMinus,
	ShieldUser,
	UserMinus,
} from "lucide-react";
import { useConfirmationDialog } from "@/hooks/use-confirmation-dialog.tsx";

type MemberDropdownMenuProps = {
	username: string;
	isAdmin: boolean;
	onChangeAdminState: () => void;
	onRemoveMember: () => void;
};

export function MemberDropdownMenu({
	username,
	isAdmin,
	onChangeAdminState,
	onRemoveMember,
}: MemberDropdownMenuProps) {
	const [openConfirmationDialog, ConfirmationDialog] = useConfirmationDialog();

	const handleRemoveMember = (username: string) => {
		openConfirmationDialog({
			title: `Are you sure you want to remove ${username}?`,
			description: `Removing ${username} will not remove any of the history associated with them, but they will no longer be able to access this session.`,
			confirmText: "Remove member",
			cancelText: "Cancel",
			onConfirm: onRemoveMember,
		});
	};

	return (
		<>
			<ConfirmationDialog />
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
						onClick={() => handleRemoveMember(username)}
					>
						<UserMinus className="size-6" />
						<span>Remove from session</span>
					</DropdownMenuItem>
				</DropdownMenuContent>
			</DropdownMenu>
		</>
	);
}
