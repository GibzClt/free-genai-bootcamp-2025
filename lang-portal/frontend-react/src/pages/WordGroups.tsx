import { useEffect, useState } from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Link } from "react-router-dom";
import { api } from "@/lib/api";
import { Loader2 } from "lucide-react";

interface WordGroup {
  id: number;
  name: string;
  description: string;
  word_count: number;
}

export default function WordGroups() {
  const [groups, setGroups] = useState<WordGroup[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchGroups = async () => {
      try {
        const response = await api.get("/api/groups");
        setGroups(response.data);
      } catch (error) {
        console.error("Failed to fetch word groups:", error);
      } finally {
        setLoading(false);
      }
    };

    fetchGroups();
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
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-3xl font-bold">Word Groups</h1>
        <Button>Create New Group</Button>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        {groups.map((group) => (
          <Link key={group.id} to={`/groups/${group.id}`}>
            <Card className="hover:bg-accent transition-colors">
              <CardHeader>
                <CardTitle>{group.name}</CardTitle>
              </CardHeader>
              <CardContent>
                <p className="text-muted-foreground mb-2">{group.description}</p>
                <p className="text-sm text-muted-foreground">
                  {group.word_count} words
                </p>
              </CardContent>
            </Card>
          </Link>
        ))}
      </div>
    </div>
  );
}
