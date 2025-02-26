
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Progress } from "@/components/ui/progress";
import { Button } from "@/components/ui/button";
import { useNavigate } from "react-router-dom";
import { useEffect, useState } from "react";

export default function Dashboard() {
  const navigate = useNavigate();
  const [studyProgress, setStudyProgress] = useState(null)
  const [quickStats, setQuickStats] = useState<any>(null);

  useEffect(() => {
    fetch(import.meta.env.VITE_API_URL + "/api/dashboard/quick-stats")
      .then((res) => res.json())
      .then((data) => setQuickStats(data));
  }, []);

  useEffect(() => {
    fetch(import.meta.env.VITE_API_URL + "/api/dashboard/study-progress")
      .then((res) => res.json())
      .then((data) => setStudyProgress(data));
  }, []);

  return (
    <div className="space-y-8">
      <div className="flex justify-between items-center">
        <h1 className="text-3xl font-bold">Dashboard</h1>
        <Button
          size="lg"
          onClick={() => navigate("/study-activities")}
          className="transition-smooth hover:scale-105"
        >
          Start Studying
        </Button>
      </div>

      <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
        <Card>
          <CardHeader>
            <CardTitle className="text-lg">Last Study Session</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-2">
              <p className="text-2xl font-semibold">Typing Tutor</p>
              <p className="text-sm text-muted-foreground">2 hours ago</p>
              <div className="mt-4 flex items-center gap-2">
                <span className="text-green-500">4 correct</span>
                <span className="text-red-500">1 wrong</span>
              </div>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle className="text-lg">Study Progress</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              <div>
                <div className="flex justify-between mb-2">
                  <span>Total Words</span>
                  <span>{`${studyProgress?.total_words_studied}/${studyProgress?.total_available_words}`}</span>
                </div>
                <Progress value={(studyProgress?.total_words_studied / studyProgress?.total_available_words) * 100} />
              </div>
              <div>
                <div className="flex justify-between mb-2">
                  <span>Mastery</span>
                  <span>{`${(studyProgress?.total_words_studied / studyProgress?.total_available_words) * 100}%`}</span>
                </div>
                <Progress value={(studyProgress?.total_words_studied / studyProgress?.total_available_words) * 100} />
              </div>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle className="text-lg">Quick Stats</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-2 gap-4">
              <div>
                <p className="text-sm text-muted-foreground">Success Rate</p>
                <p className="text-2xl font-semibold">{quickStats?.success_rate}%</p>
              </div>
              <div>
                <p className="text-sm text-muted-foreground">Study Sessions</p>
                <p className="text-2xl font-semibold">{quickStats?.total_study_sessions}</p>
              </div>
              <div>
                <p className="text-sm text-muted-foreground">Active Groups</p>
                <p className="text-2xl font-semibold">{quickStats?.total_active_groups}</p>
              </div>
              <div>
                <p className="text-sm text-muted-foreground">Study Streak</p>
                <p className="text-2xl font-semibold">{quickStats?.study_streak_days}</p>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
