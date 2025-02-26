import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
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

interface Word {
  id: number;
  japanese: string;
  reading: string;
  english: string;
}

interface WordGroup {
  id: number;
  name: string;
  description: string;
  words: Word[];
}

export default function WordGroup() {
  const { id } = useParams();
  const [group, setGroup] = useState<WordGroup | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchGroup = async () => {
      try {
        const response = await api.get(`/api/groups/${id}`);
        setGroup(response.data);
      } catch (error) {
        console.error("Failed to fetch word group:", error);
      } finally {
        setLoading(false);
      }
    };

    fetchGroup();
  }, [id]);

  if (loading) {
    return (
      <div className="flex items-center justify-center h-full">
        <Loader2 className="h-8 w-8 animate-spin" />
      </div>
    );
  }

  if (!group) {
    return <div>Group not found</div>;
  }

  return (
    <div className="container mx-auto p-6">
      <div className="flex justify-between items-center mb-6">
        <div>
          <h1 className="text-3xl font-bold mb-2">{group.name}</h1>
          <p className="text-muted-foreground">{group.description}</p>
        </div>
        <div className="space-x-2">
          <Button variant="outline">Edit Group</Button>
          <Button>Add Words</Button>
        </div>
      </div>

      <Card>
        <CardContent className="p-0">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Japanese</TableHead>
                <TableHead>Reading</TableHead>
                <TableHead>English</TableHead>
                <TableHead className="w-[100px]">Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {group.words.map((word) => (
                <TableRow key={word.id}>
                  <TableCell>{word.japanese}</TableCell>
                  <TableCell>{word.reading}</TableCell>
                  <TableCell>{word.english}</TableCell>
                  <TableCell>
                    <Button variant="ghost" size="sm">
                      Remove
                    </Button>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </CardContent>
      </Card>
    </div>
  );
}
