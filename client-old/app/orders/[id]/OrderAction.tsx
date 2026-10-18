"use client";

import { useRouter } from "next/navigation";
import { useState } from "react";

interface OrderActionsProps {
  orderID: number;
  currentStatus: string;
  hasLayout: boolean;
  isSampleFinished: boolean;
  waitingReason?: string | null;
}

export default function OrderActions({
  orderID,
  currentStatus,
  hasLayout,
  isSampleFinished,
  waitingReason,
}: OrderActionsProps) {
  const router = useRouter();
  const [loading, setLoading] = useState(false);
  const [pauseInput, setPauseInput] = useState("");
  const [showWaitingInput, setShowWaitingInput] = useState(false);
  const [sewingLine, setSewingLine] = useState("");

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
      setPauseInput("");
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
            onClick={() => executeAction("resume-order", "PATCH")}
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
            <div className="flex flex-col gap-2 p-4 border-2 border-dashed border-gray-400 rounded bg-gray-50 w-full">
              <span className="font-bold text-gray-700">
                🔒 PPIC Gate: Assign Production Line
              </span>
              <div className="flex gap-2 text-gray-900">
                <select
                  value={sewingLine}
                  onChange={(e) => setSewingLine(e.target.value)}
                  className="border p-2 rounded flex-1 bg-white"
                >
                  <option value="">-- Select Sewing Line --</option>
                  <option value="1">Line 1 (Heavy Duty / Denim)</option>
                  <option value="2">Line 2 (Delicate / Silk)</option>
                  <option value="3">Line 3 (Standard Cotton)</option>
                </select>
                <button
                  onClick={() => {
                    if (!sewingLine)
                      return alert("Please select a line first!");
                    executeAction("layout", "PATCH", {
                      layout_id: parseInt(sewingLine),
                    });
                  }}
                  className="bg-gray-800 text-white px-6 py-2 rounded hover:bg-black font-medium"
                >
                  Lock Layout
                </button>
              </div>
            </div>
          ))}

        {["IN_PRODUCTION", "CUTTING", "SEWING"].includes(currentStatus) &&
          (waitingReason ? (
            //! if pause hide next step
            <div className="flex flex-wrap items-center gap-3 border-2 border-red-500 rounded bg-red-50 w-full py-3 px-15">
              <span className="text-red-700 font-bold mr-2">
                PRODUCTION HALTED : {waitingReason}
              </span>
              <button
                onClick={() => executeAction("resume-prod", "PATCH")}
                className="bg-red-600 text-white px-4 py-2 rounded hover:bg-red-700 font-bold"
              >
                Resolve Issue
              </button>
            </div>
          ) : (
            //! if active show normal pipeline
            <>
              {currentStatus === "IN_PRODUCTION" && (
                <button
                  onClick={() => updateStatus("CUTTING")}
                  className="bg-indigo-600 text-white px-4 py-2 rounded"
                >
                  Start Cutting
                </button>
              )}
              {currentStatus === "CUTTING" && (
                <button
                  onClick={() => updateStatus("SEWING")}
                  className="bg-indigo-600 text-white px-4 py-2 rounded"
                >
                  Move to Sewing
                </button>
              )}
              {currentStatus === "SEWING" && (
                <button
                  onClick={() => updateStatus("QC")}
                  className="bg-indigo-600 text-white px-4 py-2 rounded"
                >
                  Send to QC
                </button>
              )}
              {/* The dynamic Pause button for floor managers */}
              <button
                onClick={() => setShowWaitingInput(!showWaitingInput)}
                className="bg-orange-500 text-white px-4 py-2 rounded hover:bg-orange-600 font-medium ml-auto"
              >
                Halt Production
              </button>
            </>
          ))}

        {currentStatus === "QC" && (
          <>
            <button
              onClick={() => updateStatus("FINISHING")}
              className="bg-green-600 text-white px-4 py-2 rounded"
            >
              Pass QC
            </button>
            <button
              onClick={() => updateStatus("SEWING")}
              className="bg-red-600 text-white px-4 py-2 rounded"
            >
              Fail QC (Rework)
            </button>
          </>
        )}
        {currentStatus === "FINISHING" && (
          <button
            onClick={() => updateStatus("READY_FOR_SHIPPING")}
            className="bg-blue-800 text-white px-4 py-2 rounded font-bold"
          >
            Pack & Mark Ready
          </button>
        )}

        {currentStatus === "READY_FOR_SHIPPING" && (
          <button
            onClick={() => updateStatus("SHIPPED")}
            className="bg-gray-800 text-white px-6 py-2 rounded hover:bg-black font-bold"
          >
            Mark as Shipped (Close Order)
          </button>
        )}
      </div>
      {showWaitingInput && (
        <div className="mt-4 p-4 border rounded bg-orange-50 flex gap-2">
          <input
            type="text"
            placeholder="E.g., Fabric delayed by supplier..."
            value={pauseInput}
            onChange={(e) => setPauseInput(e.target.value)}
            className="border p-2 rounded flex-1 text-gray-900"
          />
          <button
            onClick={() => {
              const isMassProduction = [
                "IN_PRODUCTION",
                "CUTTING",
                "SEWING",
                "QC",
                "FINISHING",
              ].includes(currentStatus);
              const endpoint = isMassProduction ? "hold-prod" : "waiting";

              executeAction(endpoint, "PATCH", { reason: pauseInput });
            }}
            className="bg-orange-600 text-white px-4 py-2 rounded hover:bg-orange-700 font-medium"
          >
            Confirm Pause
          </button>
        </div>
      )}
    </div>
  );
}
