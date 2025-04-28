import { Button } from "@/components/ui/button";
import {
	DropdownMenu,
	DropdownMenuContent,
	DropdownMenuGroup,
	DropdownMenuItem,
	DropdownMenuLabel,
	DropdownMenuSeparator,
	DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { useConfirmationDialog } from "@/hooks/use-confirmation-dialog";
import { ArrowLeftFromLine, CircleCheck, CircleX, EllipsisVertical } from "lucide-react";

type SessionMenuProps = {
	isActive: boolean;
	showAdminActions: boolean;
	onLeave: () => void;
	onChangeActiveState: () => void;
};
export function SessionMenu({ isActive, showAdminActions, onLeave, onChangeActiveState }: SessionMenuProps) {
	const [openConfirmationDialog, ConfirmationDialog] = useConfirmationDialog();

	const handleCloseSession = () => {
		openConfirmationDialog({
			title: `Are you sure you want to ${isActive ? "close" : "re-open"} this session?`,
			description: isActive ? closeSessionDescription : openSessionDescription,
			confirmText: `${isActive ? "Close" : "Open"} session`,
			onConfirm: onChangeActiveState,
		});
	};

	const handleLeaveSession = () => {
		openConfirmationDialog({
			title: "Are you sure you would like to leave?",
			description:
				"After leaving the session, if you would like to rejoin, you will need to be invited again.",
			confirmText: "Leave session",
			cancelText: "Stay",
			onConfirm: onLeave,
		});
	};

	return (
		<>
			<ConfirmationDialog />
			<DropdownMenu>
				<DropdownMenuTrigger asChild>
					<Button size="icon" variant="secondary" className="rounded-full">
						<EllipsisVertical />
					</Button>
				</DropdownMenuTrigger>
				<DropdownMenuContent className="min-w-64 mx-4">
					<DropdownMenuLabel className="text-lg font-semibold">Session actions</DropdownMenuLabel>
					{showAdminActions && (
						<DropdownMenuGroup>
							<DropdownMenuSeparator />
							<DropdownMenuItem className="group text-lg gap-4 cursor-pointer" onClick={handleCloseSession}>
								{isActive ? (
									<CircleX className="size-6 group-hover:text-destructive/75 transition-colors" />
								) : (
									<CircleCheck className="size-6 group-hover:text-green-500 transition-colors" />
								)}
								<span>{isActive ? "Close" : "Reopen"} session</span>
							</DropdownMenuItem>
						</DropdownMenuGroup>
					)}
					<DropdownMenuGroup>
						<DropdownMenuItem className="group text-lg gap-4 cursor-pointer" onClick={handleLeaveSession}>
							<ArrowLeftFromLine className="size-6 text-muted-foreground group-hover:text-primary/75 transition-colors" />
							<span>Leave session</span>
						</DropdownMenuItem>
					</DropdownMenuGroup>
				</DropdownMenuContent>
			</DropdownMenu>
		</>
	);
}

const closeSessionDescription = (
	<div className="flex flex-col gap-2">
		<p>
			Once closed, members will be able to see the session, but will not be able to interact with it until it
			is reopened.
		</p>
		<p>The session can be reopened by any admin member.</p>
	</div>
);

const openSessionDescription = "Once re-opened, members will be able to interact with this session again.";

SessionMenu.Skeleton = function SessionMenuSkeleton() {
	return (
		<Button size="icon" variant="secondary" className="rounded-full animate-pulse">
			<EllipsisVertical />
		</Button>
	);
};
