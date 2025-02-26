
import { useState } from "react";
import { Plus, Search, Pencil, Trash2 } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from "@/components/ui/dialog";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { toast } from "sonner";

type Word = {
  id: string;
  original: string;
  translation: string;
  group: string;
  lastReviewed: string;
};

// Temporary mock data - will be replaced with API call
const mockWords: Word[] = [
  {
    id: "1",
    original: "Hello",
    translation: "Hola",
    group: "Basic Greetings",
    lastReviewed: "2024-03-10",
  },
  {
    id: "2",
    original: "Goodbye",
    translation: "Adiós",
    group: "Basic Greetings",
    lastReviewed: "2024-03-09",
  },
  {
    id: "3",
    original: "Thank you",
    translation: "Gracias",
    group: "Common Phrases",
    lastReviewed: "2024-03-08",
  },
];

const mockGroups = [
  "Basic Greetings",
  "Common Phrases",
  "Numbers",
  "Colors",
  "Food and Drinks",
];

export default function Words() {
  const [searchTerm, setSearchTerm] = useState("");
  const [words, setWords] = useState<Word[]>(mockWords);
  const [isDialogOpen, setIsDialogOpen] = useState(false);
  const [editingWord, setEditingWord] = useState<Word | null>(null);
  const [formData, setFormData] = useState({
    original: "",
    translation: "",
    group: "",
  });

  const filteredWords = words.filter(
    (word) =>
      word.original.toLowerCase().includes(searchTerm.toLowerCase()) ||
      word.translation.toLowerCase().includes(searchTerm.toLowerCase()) ||
      word.group.toLowerCase().includes(searchTerm.toLowerCase())
  );

  const handleOpenDialog = (word?: Word) => {
    if (word) {
      setEditingWord(word);
      setFormData({
        original: word.original,
        translation: word.translation,
        group: word.group,
      });
    } else {
      setEditingWord(null);
      setFormData({
        original: "",
        translation: "",
        group: "",
      });
    }
    setIsDialogOpen(true);
  };

  const handleCloseDialog = () => {
    setIsDialogOpen(false);
    setEditingWord(null);
    setFormData({
      original: "",
      translation: "",
      group: "",
    });
  };

  const handleSave = () => {
    if (!formData.original || !formData.translation || !formData.group) {
      toast.error("Please fill in all fields");
      return;
    }

    if (editingWord) {
      setWords(words.map(word => 
        word.id === editingWord.id 
          ? { ...word, ...formData }
          : word
      ));
      toast.success("Word updated successfully");
    } else {
      const newWord: Word = {
        id: Math.random().toString(36).substr(2, 9),
        ...formData,
        lastReviewed: new Date().toISOString(),
      };
      setWords([...words, newWord]);
      toast.success("Word added successfully");
    }
    handleCloseDialog();
  };

  const handleDelete = (id: string) => {
    setWords(words.filter(word => word.id !== id));
    toast.success("Word deleted successfully");
  };

  return (
    <div className="space-y-8">
      <div className="flex items-center justify-between">
        <h1 className="text-3xl font-bold">Words</h1>
        <Button size="lg" onClick={() => handleOpenDialog()}>
          <Plus className="h-4 w-4 mr-2" />
          Add Word
        </Button>
      </div>

      <div className="flex items-center space-x-2">
        <div className="relative flex-1">
          <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
          <Input
            placeholder="Search words..."
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            className="pl-9"
          />
        </div>
      </div>

      <div className="rounded-md border">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Original</TableHead>
              <TableHead>Translation</TableHead>
              <TableHead>Group</TableHead>
              <TableHead>Last Reviewed</TableHead>
              <TableHead className="w-[100px]">Actions</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {filteredWords.map((word) => (
              <TableRow key={word.id}>
                <TableCell className="font-medium">{word.original}</TableCell>
                <TableCell>{word.translation}</TableCell>
                <TableCell>{word.group}</TableCell>
                <TableCell>{new Date(word.lastReviewed).toLocaleDateString()}</TableCell>
                <TableCell>
                  <div className="flex items-center gap-2">
                    <Button 
                      variant="ghost" 
                      size="icon"
                      onClick={() => handleOpenDialog(word)}
                    >
                      <Pencil className="h-4 w-4" />
                    </Button>
                    <Button 
                      variant="ghost" 
                      size="icon"
                      onClick={() => handleDelete(word.id)}
                    >
                      <Trash2 className="h-4 w-4" />
                    </Button>
                  </div>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </div>

      <Dialog open={isDialogOpen} onOpenChange={setIsDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>
              {editingWord ? "Edit Word" : "Add New Word"}
            </DialogTitle>
          </DialogHeader>
          <div className="space-y-4 py-4">
            <div className="space-y-2">
              <label className="text-sm font-medium">Original Word</label>
              <Input
                value={formData.original}
                onChange={(e) =>
                  setFormData({ ...formData, original: e.target.value })
                }
                placeholder="Enter original word"
              />
            </div>
            <div className="space-y-2">
              <label className="text-sm font-medium">Translation</label>
              <Input
                value={formData.translation}
                onChange={(e) =>
                  setFormData({ ...formData, translation: e.target.value })
                }
                placeholder="Enter translation"
              />
            </div>
            <div className="space-y-2">
              <label className="text-sm font-medium">Group</label>
              <Select
                value={formData.group}
                onValueChange={(value) =>
                  setFormData({ ...formData, group: value })
                }
              >
                <SelectTrigger>
                  <SelectValue placeholder="Select a group" />
                </SelectTrigger>
                <SelectContent>
                  {mockGroups.map((group) => (
                    <SelectItem key={group} value={group}>
                      {group}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={handleCloseDialog}>
              Cancel
            </Button>
            <Button onClick={handleSave}>
              {editingWord ? "Save Changes" : "Add Word"}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}
