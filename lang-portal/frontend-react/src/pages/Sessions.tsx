import { useEffect, useState } from "react";
import { Card, CardContent } from "@/components/ui/card";
import { api } from "@/lib/api";
import { Loader2 } from "lucide-react";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Link } from "react-router-dom";
import { format } from "date-fns";

interface StudySession {
  id: number;
  activity_name: string;
  group_name: string;
  start_time: string;
  end_time: string;
  review_items_count: number;
}

export default function Sessions() {
  const [sessions, setSessions] = useState<StudySession[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchSessions = async () => {
      try {
        const response = await api.get("/api/study-sessions");
        setSessions(response.data);
      } catch (error) {
        console.error("Failed to fetch study sessions:", error);
      } finally {
        setLoading(false);
      }
    };

    fetchSessions();
  }, []);

  if (loading) {
    return (
      <div className="flex items-center justify-center h-full">
        <Loader2 className="h-8 w-8 animate-spin" />
      </div>
    );
  }

  return (
    <div className="container mx-auto p-6">
      <h1 className="text-3xl font-bold mb-6">Study Sessions</h1>

      <Card>
        <CardContent className="p-0">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Activity</TableHead>
                <TableHead>Group</TableHead>
                <TableHead>Start Time</TableHead>
                <TableHead>End Time</TableHead>
                <TableHead>Items Reviewed</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {sessions.map((session) => (
                <TableRow key={session.id}>
                  <TableCell>
                    <Link
                      to={`/study-activities/${session.id}`}
                      className="text-primary hover:underline"
                    >
                      {session.activity_name}
                    </Link>
                  </TableCell>
                  <TableCell>{session.group_name}</TableCell>
                  <TableCell>
                    {format(new Date(session.start_time), "PPp")}
                  </TableCell>
                  <TableCell>
                    {session.end_time
                      ? format(new Date(session.end_time), "PPp")
                      : "In Progress"}
                  </TableCell>
                  <TableCell>{session.review_items_count}</TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </CardContent>
      </Card>
    </div>
  );
}
