
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Play, Eye } from "lucide-react";

type StudyActivity = {
  id: string;
  name: string;
  thumbnail: string;
};

// Temporary mock data - will be replaced with API call
const activities: StudyActivity[] = [
  {
    id: "1",
    name: "Typing Tutor",
    thumbnail: "/placeholder.svg",
  },
  {
    id: "2",
    name: "Flashcards",
    thumbnail: "/placeholder.svg",
  },
  {
    id: "3",
    name: "Multiple Choice",
    thumbnail: "/placeholder.svg",
  },
];

export default function StudyActivities() {
  return (
    <div className="space-y-8">
      <div className="flex items-center justify-between">
        <h1 className="text-3xl font-bold">Study Activities</h1>
      </div>

      <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
        {activities.map((activity) => (
          <Card key={activity.id} className="overflow-hidden">
            <CardHeader className="p-0">
              <div className="aspect-video w-full overflow-hidden bg-muted">
                <img
                  src={activity.thumbnail}
                  alt={activity.name}
                  className="h-full w-full object-cover"
                />
              </div>
            </CardHeader>
            <CardContent className="p-6">
              <CardTitle className="text-xl">{activity.name}</CardTitle>
            </CardContent>
            <CardFooter className="px-6 pb-6 pt-0 gap-3">
              <Button 
                className="flex-1"
                onClick={() => window.location.href = `/study-activities/${activity.id}/launch`}
              >
                <Play className="h-4 w-4 mr-2" />
                Launch
              </Button>
              <Button 
                variant="outline" 
                className="flex-1"
                onClick={() => window.location.href = `/study-activities/${activity.id}`}
              >
                <Eye className="h-4 w-4 mr-2" />
                View
              </Button>
            </CardFooter>
          </Card>
        ))}
      </div>
    </div>
  );
}
