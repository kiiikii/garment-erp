"use client";

import { useRouter } from "next/navigation";
import { useState } from "react";

interface OrderActionsProps {
  orderID: number;
  currentStatus: string;
  hasLayout: boolean;
  isSampleFinished: boolean;
}

export default function OrderActions({
  orderID,
  currentStatus,
  hasLayout,
  isSampleFinished,
}: OrderActionsProps) {
  const router = useRouter();
  const [loading, setLoading] = useState(false);
  const [waitingReason, setWaitingReason] = useState("");
  const [showWaitingInput, setShowWaitingInput] = useState(false);

  //! helper for standard status updates
  const updateStatus = async (newStatus: string) => {
    setLoading(true);
    try {
      const res = await fetch(
        `http://localhost:8080/orders/status?id=${orderID}`,
        {
          method: "PATCH",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({
            current_status: currentStatus,
            new_status: newStatus,
          }),
        },
      );
      if (!res.ok) throw new Error(await res.text());
      router.refresh(); // Magically updates the Next.js server page!
    } catch (err: unknown) {
      if (err instanceof Error) {
        alert("Error: " + err.message);
      } else {
        alert("An unknown error occurred");
      }
    } finally {
      setLoading(false);
    }
  };

  //! helper custom endpoint
  const executeAction = async (
    endpoint: string,
    method: string,
    body?: Record<string, unknown>,
  ) => {
    setLoading(true);
    try {
      const res = await fetch(
        `http://localhost:8080/orders/${endpoint}?id=${orderID}`,
        {
          method,
          headers: body ? { "Content-Type": "application/json" } : undefined,
          body: body ? JSON.stringify(body) : undefined,
        },
      );
      if (!res.ok) throw new Error(await res.text());
      setShowWaitingInput(false);
      setWaitingReason("");
      router.refresh();
    } catch (err: unknown) {
      if (err instanceof Error) {
        alert("Error: " + err.message);
      } else {
        alert("An unknown error occurred");
      }
    } finally {
      setLoading(false);
    }
  };

  if (loading)
    return (
      <div className="p-4 bg-gray-100 rounded text-gray-500 font-bold">
        Processing Action...
      </div>
    );

  return (
    <div className="bg-white p-6 rounded-lg shadow space-y-4 border-l-4 border-blue-500">
      <h2 className="text-lg font-bold text-gray-800 border-b pb-2">
        Available Actions
      </h2>
      <div className="flex flex-wrap gap-3">
        {/* sampling action */}
        {currentStatus === "WAITING_FOR_SAMPLE" && !isSampleFinished && (
          <button
            onClick={() => executeAction("sample-complete", "PATCH")}
            className="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700 font-medium"
          >
            Mark sample finished
          </button>
        )}

        {currentStatus === "WAITING_FOR_SAMPLE" && isSampleFinished && (
          <>
            <button
              onClick={() => updateStatus("SAMPLE_APPROVED")}
              className="bg-green-600 text-white px-4 py-2 rounded hover:bg-green-700 font-medium"
            >
              Approved Sample
            </button>
            <button
              onClick={() => updateStatus("SAMPLE_REVISED")}
              className="bg-yellow-500 text-white px-4 py-2 rounded hover:bg-yellow-600 font-medium"
            >
              Revise Sample
            </button>
            <button
              onClick={() => updateStatus("SAMPLE_REJECTED")}
              className="bg-red-600 text-white px-4 py-2 rounded hover:bg-red-700 font-medium"
            >
              Reject Sample
            </button>
          </>
        )}

        {currentStatus === "SAMPLE_REVISED" && (
          <button
            onClick={() => updateStatus("WAITING_FOR_SAMPLE")}
            className="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700 font-medium"
          >
            Acknowledge & Start Rework (Back to Waiting)
          </button>
        )}

        {currentStatus === "SAMPLE_APPROVED" && (
          <>
            <button
              onClick={() => updateStatus("LOA_SIGNED")}
              className="bg-purple-600 text-white px-4 py-2 rounded hover:bg-purple-700 font-medium"
            >
              Sign LoA
            </button>
            <button
              onClick={() => updateStatus("LOA_REVISED")}
              className="bg-yellow-500 text-white px-4 py-2 rounded hover:bg-yellow-600 font-medium"
            >
              Revise LoA Terms
            </button>
            <button
              onClick={() => updateStatus("LOA_REJECTED")}
              className="bg-red-600 text-white px-4 py-2 rounded hover:bg-red-700 font-medium"
            >
              Reject Order
            </button>
          </>
        )}

        {currentStatus === "LOA_REVISED" && (
          <div className="flex flex-wrap items-center gap-3 border p-3 rounded bg-yellow-50 w-full">
            <span className="text-yellow-700 font-bold mr-2">
              Negotiating new terms...
            </span>
            <button
              onClick={() => updateStatus("SAMPLE_APPROVED")}
              className="bg-purple-600 text-white px-4 py-2 rounded hover:bg-purple-700 font-medium"
            >
              Submit New LoA Terms (Back to Waiting)
            </button>
          </div>
        )}

        {currentStatus === "LOA_SIGNED" && (
          <>
            <button
              onClick={() => updateStatus("READY_FOR_PRODUCTION")}
              className="bg-green-600 text-white px-4 py-2 rounded hover:bg-green-700 font-medium"
            >
              Mark Ready (Materials Arrived)
            </button>
            <button
              onClick={() => setShowWaitingInput(!showWaitingInput)}
              className="bg-orange-500 text-white px-4 py-2 rounded hover:bg-orange-600 font-medium"
            >
              Pause for Materials
            </button>
          </>
        )}

        {currentStatus === "WAITING_FOR_MATERIALS" && (
          <button
            onClick={() => executeAction("resume", "PATCH")}
            className="bg-green-600 text-white px-4 py-2 rounded hover:bg-green-700 font-medium"
          >
            Resume Order (Materials Arrived)
          </button>
        )}

        {currentStatus === "READY_FOR_PRODUCTION" &&
          (hasLayout ? (
            <button
              onClick={() => updateStatus("IN_PRODUCTION")}
              className="bg-blue-800 text-white px-4 py-2 rounded hover:bg-blue-900 font-bold"
            >
              Send to Production Floor
            </button>
          ) : (
            <span className="text-red-600 font-bold px-2 py-2">
              Cannot start production: Assign PPIC Layout first.
            </span>
          ))}
      </div>
      {showWaitingInput && (
        <div className="mt-4 p-4 border rounded bg-orange-50 flex gap-2">
          <input
            type="text"
            placeholder="E.g., Fabric delayed by supplier..."
            value={waitingReason}
            onChange={(e) => setWaitingReason(e.target.value)}
            className="border p-2 rounded flex-1"
          />
          <button
            onClick={() =>
              executeAction("waiting", "PATCH", { reason: waitingReason })
            }
            className="bg-orange-600 text-white px-4 py-2 rounded hover:bg-orange-700 font-medium"
          >
            Confirm Pause
          </button>
        </div>
      )}
    </div>
  );
}
