import {
	AlertDialog,
	AlertDialogAction,
	AlertDialogContent,
	AlertDialogDescription,
	AlertDialogFooter,
	AlertDialogHeader,
	AlertDialogTitle,
} from "@/components/ui/alert-dialog";
import { Button } from "@/components/ui/button.tsx";
import { type JSX, useCallback, useState } from "react";

interface ConfirmationDialogOptions {
	title: string;
	description?: string;
	confirmText?: string;
}

export function useInformationDialog(): [
	(options: ConfirmationDialogOptions) => void,
	() => JSX.Element | null,
] {
	const [open, setOpen] = useState(false);
	const [options, setOptions] = useState<ConfirmationDialogOptions | null>(
		null,
	);

	const openDialog = useCallback((opts: ConfirmationDialogOptions) => {
		setOptions(opts);
		setOpen(true);
	}, []);

	const handleCancel = useCallback(() => {
		setOpen(false);
	}, []);

	const InformationDialog = () => {
		if (!options) return null;
		return (
			<AlertDialog open={open} onOpenChange={setOpen}>
				<AlertDialogContent>
					<AlertDialogHeader>
						<AlertDialogTitle className="text-left text-2xl">
							{options.title}
						</AlertDialogTitle>
						{options.description && (
							<AlertDialogDescription className="text-left text-lg">
								{options.description}
							</AlertDialogDescription>
						)}
					</AlertDialogHeader>
					<AlertDialogFooter>
						<section className="flex gap-2 justify-end tracking-wider">
							<AlertDialogAction onClick={handleCancel} asChild>
								<Button variant="secondary" size="lg">
									{options.confirmText ?? "OK"}
								</Button>
							</AlertDialogAction>
						</section>
					</AlertDialogFooter>
				</AlertDialogContent>
			</AlertDialog>
		);
	};

	return [openDialog, InformationDialog];
}
