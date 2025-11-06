type FeedbackType = "success" | "info" | "warning" | "error";

type Feedback = {
    type: FeedbackType;
    message: string;
}

export type {Feedback, FeedbackType}