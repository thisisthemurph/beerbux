import { Card, CardContent } from "@/components/ui/card";
import { useInformationDialog } from "@/hooks/use-information-dialog.tsx";
import { MinusCircle, PlusCircle } from "lucide-react";

type UserCardProps = {
	username: string;
	netBalance: number;
};

export function UserCard({ username, netBalance }: UserCardProps) {
	return (
		<Card className="bg-primary text-white min-h-36 shadow-xl">
			<CardContent className="flex flex-col justify-between">
				<section className="flex justify-between font-semibold">
					<div>
						<p className="text-2xl tracking-wider">Beerbux</p>
						<p className="text-sm font-mono tracking-wide">{username}</p>
					</div>
					<NetBalance balance={netBalance} />
				</section>
			</CardContent>
		</Card>
	);
}

function NetBalance({ balance }: { balance: number }) {
	const [openInformationDialog, InformationDialog] = useInformationDialog();

	const handleInformationClick = () => {
		openInformationDialog({
			title: "Net Balance",
			description: (
				<>
					<p>Think of this number as your beer bank balance of all beers bought and received.</p>
					<div className="grid grid-cols-5 gap-4 mt-4">
						<PlusCircle className="size-8 col-span-1 text-green-500 translate-y-1/2" />
						<p className="col-span-4">A positive number means you received more than you bought.</p>
						<MinusCircle className="size-8 col-span-1 text-red-500 translate-y-1/2" />
						<p className="col-span-4">A negative number means you're either a stand-up guy or a pushover.</p>
					</div>
				</>
			),
		});
	};

	return (
		<>
			<InformationDialog />
			<button type="button" title="Total of all beers given and recieved" onClick={handleInformationClick}>
				<p className="text-right">
					<span className="text-lg">$</span>
					<span className="text-3xl text-right tracking-wide">{balance ?? 0}</span>
				</p>
				<p className="text-xs text-right tracking-wide font-semibold">net</p>
			</button>
		</>
	);
}
