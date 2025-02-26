
import { useParams } from "react-router-dom";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { useState } from "react";

// Temporary mock data - will be replaced with API call
const mockGroups = [
  { id: "1", name: "Basic Verbs" },
  { id: "2", name: "Common Nouns" },
  { id: "3", name: "Adjectives" },
];

export default function StudyActivityLauncher() {
  const { id } = useParams();
  const [selectedGroup, setSelectedGroup] = useState<string>("");
  const activity = {
    id,
    name: "Typing Tutor",
  };

  const handleLaunch = () => {
    // This will be replaced with actual API call
    window.open(`/study/${activity.id}?group=${selectedGroup}`, '_blank');
    window.location.href = `/study-activities/${id}`;
  };

  return (
    <div className="max-w-md mx-auto space-y-8 pt-8">
      <h1 className="text-3xl font-bold text-center">{activity.name}</h1>

      <Card>
        <CardHeader>
          <CardTitle>Launch Activity</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="space-y-2">
            <label className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70">
              Select Word Group
            </label>
            <Select
              value={selectedGroup}
              onValueChange={setSelectedGroup}
            >
              <SelectTrigger>
                <SelectValue placeholder="Select a group" />
              </SelectTrigger>
              <SelectContent>
                {mockGroups.map((group) => (
                  <SelectItem key={group.id} value={group.id}>
                    {group.name}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>

          <Button 
            className="w-full" 
            size="lg"
            onClick={handleLaunch}
            disabled={!selectedGroup}
          >
            Launch Now
          </Button>
        </CardContent>
      </Card>
    </div>
  );
}
