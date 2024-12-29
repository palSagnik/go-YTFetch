export interface Video {
    id: string;
    title: string;
    description: string;
    published_at: string;
    channeltitle: string;
    thumbnailurl: string;
  }
  
  export interface ApiResponse {
    status: string;
    data: Video[];
    pagination: {
      hasNext: boolean;
      nextCursor: string;
      totalCount: number;
    };
  }