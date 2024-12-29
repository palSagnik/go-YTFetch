// src/components/VideoList.tsx
import { useEffect, useState } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { Card, CardContent, CardFooter } from "./ui/card";
import { Button } from "./ui/button";
import { ChevronLeft, ChevronRight, Calendar, User, Play } from "lucide-react";
import { Video } from "../types";

interface ApiResponse {
  status: string;
  data: Video[];
  pagination: {
    hasNext: boolean;
    nextCursor: string;
    totalCount: number;
  };
}

const VideoList = () => {
  const [videos, setVideos] = useState<Video[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [pagination, setPagination] = useState<ApiResponse["pagination"] | null>(null);
  const [hoveredVideo, setHoveredVideo] = useState<string | null>(null);
  
  const { cursor } = useParams();
  const navigate = useNavigate();

  useEffect(() => {
    const fetchVideos = async () => {
      try {
        setIsLoading(true);
        const url = new URL("http://localhost:9000/api/getVideos");
        if (cursor) {
          const decodedCursor = decodeURIComponent(cursor);
          url.searchParams.append("next_cursor", decodedCursor);
        }

        const response = await fetch(url.toString());
        if (!response.ok) {
          throw new Error("Failed to fetch videos");
        }

        const data: ApiResponse = await response.json();
        console.log("API Response:", data);
        setVideos(data.data);
        setPagination(data.pagination);
      } catch (err) {
        console.error("Fetch error:", err);
        setError(err instanceof Error ? err.message : "An error occurred");
      } finally {
        setIsLoading(false);
      }
    };

    fetchVideos();
  }, [cursor]);

  const goToNextPage = () => {
    if (pagination?.hasNext && pagination?.nextCursor) {
      const encodedCursor = encodeURIComponent(pagination.nextCursor);
      navigate(`/${encodedCursor}`);
    }
  };

  const goToPreviousPage = () => {
    navigate("/");
  };

  if (isLoading) {
    return (
      <div className="max-w-7xl mx-auto px-4 py-8 space-y-12 bg-gray-900 min-h-screen flex items-center justify-center">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-purple-500"></div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="max-w-7xl mx-auto px-4 py-8 space-y-12 bg-gray-900 min-h-screen flex items-center justify-center">
        <div className="text-red-500 text-center">Error: {error}</div>
      </div>
    );
  }

  return (
    <div className="max-w-7xl mx-auto px-4 py-8 space-y-12 bg-gray-900">
      <div className="relative overflow-hidden rounded-2xl bg-gradient-to-r from-indigo-900 to-purple-900 p-8 text-white">
        <div className="relative z-10">
          <h1 className="text-4xl font-bold mb-4">Featured Videos</h1>
          <p className="text-lg text-gray-300">Discover trending content</p>
        </div>
        <div className="absolute top-0 right-0 w-64 h-64 bg-white opacity-5 rounded-full -translate-y-1/2 translate-x-1/2" />
        <div className="absolute bottom-0 left-0 w-32 h-32 bg-white opacity-5 rounded-full translate-y-1/2 -translate-x-1/2" />
      </div>

      <ul className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-8">
        {videos.map((video) => (
          <li
            key={video.id}
            className="group"
            onMouseEnter={() => setHoveredVideo(video.id)}
            onMouseLeave={() => setHoveredVideo(null)}
          >
            <Card className="h-full overflow-hidden transition-all duration-500 hover:shadow-2xl bg-gray-800 border-gray-700 hover:border-purple-500 flex flex-col">
              <CardContent className="p-0">
                <div className="relative w-full pt-[56.25%] overflow-hidden">
                  <img
                    src={video.thumbnailurl
                      .replace("/default", "/hqdefault")
                      .replace("/default_live", "/hqdefault_live")}
                    alt={video.title}
                    className="absolute inset-0 w-full h-full object-cover transform transition-transform duration-700 group-hover:scale-110"
                  />
                  <div
                    className={`absolute inset-0 bg-gradient-to-t from-black/80 to-transparent transition-opacity duration-300 ${
                      hoveredVideo === video.id ? "opacity-100" : "opacity-0"
                    }`}
                  />
                  <div
                    className={`absolute inset-0 flex items-center justify-center transition-opacity duration-300 ${
                      hoveredVideo === video.id ? "opacity-100" : "opacity-0"
                    }`}
                  >
                    <div className="w-16 h-16 bg-purple-500/20 backdrop-blur-sm rounded-full flex items-center justify-center transform transition-transform duration-300 hover:scale-110">
                      <Play className="h-8 w-8 text-white fill-current" />
                    </div>
                  </div>
                </div>
                <div className="p-6 space-y-4">
                  <h2 className="text-lg font-semibold line-clamp-2 h-14 text-gray-100 group-hover:text-purple-400 transition-colors duration-300">
                    {video.title}
                  </h2>
                  <div className="flex items-center space-x-2 text-gray-400">
                    <User className="h-4 w-4" />
                    <p className="text-sm font-medium">{video.channeltitle}</p>
                  </div>
                  <p className="text-sm text-gray-500 line-clamp-2">
                    {video.description}
                  </p>
                </div>
              </CardContent>
              <CardFooter className="bg-gray-800/50 border-t border-gray-700">
                <div className="flex items-center justify-center space-x-2 text-gray-400">
                  <Calendar className="h-4 w-4" />
                  <p className="text-sm">
                    {new Date(video.published_at).toLocaleDateString(
                      undefined,
                      {
                        year: "numeric",
                        month: "long",
                        day: "numeric",
                      }
                    )}
                  </p>
                </div>
              </CardFooter>
            </Card>
          </li>
        ))}
      </ul>

      <div className="flex justify-between items-center pt-8 border-t border-gray-700">
        <div className="flex items-center space-x-4">
          <Button
            onClick={goToPreviousPage}
            disabled={!cursor}
            variant="outline"
            className="flex items-center space-x-2 px-6 py-2 transition-all duration-300 bg-gray-800 border-gray-700 hover:bg-gray-700 hover:text-purple-400 hover:border-purple-500 disabled:opacity-50 text-gray-300"
          >
            <ChevronLeft className="h-4 w-4" />
            <span>Previous</span>
          </Button>
          <Button
            onClick={goToNextPage}
            disabled={!pagination?.hasNext}
            variant="outline"
            className="flex items-center space-x-2 px-6 py-2 transition-all duration-300 bg-gray-800 border-gray-700 hover:bg-gray-700 hover:text-purple-400 hover:border-purple-500 disabled:opacity-50 text-gray-300"
          >
            <span>Next</span>
            <ChevronRight className="h-4 w-4" />
          </Button>
        </div>

        <div className="px-4 py-2 bg-gray-800 text-purple-400 border border-purple-500/30 rounded-full text-sm font-medium">
          {pagination?.totalCount || 0} videos
        </div>
      </div>
    </div>
  );
};

export default VideoList;