import { Card, CardContent } from "@/components/ui/card";

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
					<div title="Total of all beers given and recieved">
						<p className="text-right">
							<span className="text-lg">$</span>
							<span className="text-3xl text-right tracking-wide">
								{netBalance ?? 0}
							</span>
						</p>
						<p className="text-xs text-right tracking-wide font-semibold">
							net
						</p>
					</div>
				</section>
			</CardContent>
		</Card>
	);
}
