
import { useParams } from "react-router-dom";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Play } from "lucide-react";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";

type StudySession = {
  id: string;
  activityName: string;
  groupName: string;
  startTime: string;
  endTime: string;
  reviewItems: number;
};

// Temporary mock data - will be replaced with API call
const mockSessions: StudySession[] = [
  {
    id: "1",
    activityName: "Typing Tutor",
    groupName: "Basic Verbs",
    startTime: "2024-03-10T10:00:00",
    endTime: "2024-03-10T10:15:00",
    reviewItems: 20,
  },
  {
    id: "2",
    activityName: "Typing Tutor",
    groupName: "Common Nouns",
    startTime: "2024-03-09T15:30:00",
    endTime: "2024-03-09T15:45:00",
    reviewItems: 15,
  },
];

export default function StudyActivityDetails() {
  const { id } = useParams();
  const activity = {
    id,
    name: "Typing Tutor",
    thumbnail: "/placeholder.svg",
  };

  return (
    <div className="space-y-8">
      <div className="flex items-center justify-between">
        <h1 className="text-3xl font-bold">{activity.name}</h1>
        <Button
          size="lg"
          onClick={() => window.location.href = `/study-activities/${id}/launch`}
        >
          <Play className="h-4 w-4 mr-2" />
          Launch Activity
        </Button>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Activity Preview</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="aspect-video w-full overflow-hidden rounded-lg bg-muted">
            <img
              src={activity.thumbnail}
              alt={activity.name}
              className="h-full w-full object-cover"
            />
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>Study Sessions</CardTitle>
        </CardHeader>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>ID</TableHead>
                <TableHead>Activity</TableHead>
                <TableHead>Group</TableHead>
                <TableHead>Start Time</TableHead>
                <TableHead>End Time</TableHead>
                <TableHead className="text-right">Review Items</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {mockSessions.map((session) => (
                <TableRow key={session.id}>
                  <TableCell className="font-medium">{session.id}</TableCell>
                  <TableCell>{session.activityName}</TableCell>
                  <TableCell>{session.groupName}</TableCell>
                  <TableCell>{new Date(session.startTime).toLocaleString()}</TableCell>
                  <TableCell>{new Date(session.endTime).toLocaleString()}</TableCell>
                  <TableCell className="text-right">{session.reviewItems}</TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </CardContent>
      </Card>
    </div>
  );
}
