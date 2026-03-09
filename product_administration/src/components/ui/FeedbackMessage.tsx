interface FeedbackProps {
    type: 'error' | 'success';
    message: string;
}

export default function FeedbackMessage({ type, message }: FeedbackProps) {
    const baseClasses = "p-4 text-sm rounded-lg border";
    
    const colorClasses = type === 'error'
        ? "text-red-800 bg-red-50 border-red-100"
        : "text-green-800 bg-green-50 border-green-100";

    return (
        <div className={`${baseClasses} ${colorClasses}`} role="alert">
            {type === 'error' && <span className="font-medium">Hiba: </span>}
            {message}
        </div>
    );
}