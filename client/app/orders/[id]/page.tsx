// export const dynamic = "force-dynamic";

import Image from "next/image";
import OrderActions from "./OrderAction";
import Link from "next/link";

interface OrderSize {
  size_label: string;
  quantity: number;
}

interface OrderImage {
  id: number;
  image_url: string;
}

interface Order {
  id: number;
  total_quantity: number;
  production_type: string;
  status: string;
  created_at: string;
  internal_sample_deadline: string;
  customer_sample_deadline: string;
  actual_sample_finished_at: string | null;
  layout_id: number | null;
  customer: {
    id: number;
    name: string;
    phone: string;
    address: string;
  };
  sizes: OrderSize[] | null;
  images: OrderImage[] | null;
  waiting_reason?: string | null;
  internal_production_deadline: string;
  customer_production_deadline: string;
  actual_production_started_at: string | null;
  actual_production_completed_at: string | null;
}

interface PageProps {
  params: Promise<{ id: string }>;
}

export default async function OrderDetailPage({ params }: PageProps) {
  const resolvedParams = await params;
  const orderID = resolvedParams.id;

  const res = await fetch(`http://localhost:8080/orders?id=${orderID}`, {
    cache: "no-cache",
  });

  if (!res.ok) {
    return (
      <main className="min-h-screen bg-gray-50 p-8 font-sans">
        <div className="max-w-3xl mx-auto bg-white p-6 rounded-lg shadow text-center">
          <h1 className="text-2xl font-bold text-red-600 mb-2">
            Order Not Found
          </h1>
          <p className="text-gray-600 mb-4">
            Could not retrieve order #{orderID} from the backend.
          </p>
          <Link
            href="/orders"
            className="text-blue-600 hover:underline font-medium"
          >
            ← Back to Orders Dashboard
          </Link>
        </div>
      </main>
    );
  }

  const getCurrentLocation = (status: string, layoutID: number | null) => {
    if (!layoutID) {
      return (
        <span className="text-red-600">Unassigned (PPIC Gate Locked)</span>
      );
    }

    switch (status) {
      case "CUTTING":
        return <span className="text-indigo-600 font-bold">Cutting Room</span>;
      case "SEWING":
        return (
          <span className="text-blue-600 font-bold">
            Sewing Line #{layoutID}
          </span>
        );
      case "QC":
        return (
          <span className="text-orange-600 font-bold">
            Quality Control (QC) Station
          </span>
        );
      case "FINISHING":
        return (
          <span className="text-green-600 font-bold">Finishing & Packing</span>
        );
      case "READY_FOR_SHIPPING":
        return (
          <span className="text-gray-800 font-bold">
            Warehouse (Ready to Ship)
          </span>
        );
      case "SHIPPED":
        return (
          <span className="text-black font-extrabold">Shipped & Closed</span>
        );
      default:
        return (
          <span className="text-gray-600">
            Assigned to Sewing Line #{layoutID}
          </span>
        );
    }
  };

  const orders: Order = await res.json();

  return (
    <main className="min-h-screen bg-gray-50 p-8 font-sans">
      <div className="max-w-4xl mx-auto space-y-6">
        {/* Navigation & Header */}
        <div className="flex justify-between items-center">
          <div>
            <Link
              href="/orders"
              className="text-blue-600 hover:underline text-sm font-medium"
            >
              ← Back to All Orders
            </Link>
            <h1 className="text-3xl font-bold text-gray-800 mt-1">
              Order #{orders.id} Details
            </h1>
          </div>
          {orders.waiting_reason ? (
            <span className="bg-red-600 text-white px-4 py-2 rounded-full font-bold shadow animate-pulse">
              HALTED ({orders.status})
            </span>
          ) : (
            <span
              className={`px-3 py-1.5 rounded-full text-sm font-bold ${
                orders.status === "IN_PRODUCTION"
                  ? "bg-green-100 text-green-800"
                  : orders.status === "LOA_SIGNED"
                    ? "bg-purple-100 text-purple-800"
                    : orders.status === "SAMPLE_APPROVED"
                      ? "bg-blue-100 text-blue-800"
                      : "bg-yellow-100 text-yellow-800"
              }`}
            >
              {orders.status}
            </span>
          )}
        </div>

        {/* Order Action */}
        <OrderActions
          orderID={orders.id}
          currentStatus={orders.status}
          hasLayout={orders.layout_id !== null}
          isSampleFinished={orders.actual_sample_finished_at !== null}
          waitingReason={orders.waiting_reason}
        />

        {/* Grid Overview */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          {/* Customer & Production Info */}
          <div className="bg-white p-6 rounded-lg shadow space-y-4">
            <h2 className="text-lg font-bold text-gray-800 border-b pb-2">
              Customer & Spec
            </h2>
            <div>
              <p className="text-xs text-gray-500 uppercase font-semibold">
                Customer Name
              </p>
              <p className="text-gray-800 font-medium">
                {orders.customer?.name}
              </p>
            </div>
            <div>
              <p className="text-xs text-gray-500 uppercase font-semibold">
                Production Type
              </p>
              <p className="text-blue-600 font-bold">
                {orders.production_type}
              </p>
            </div>
            <div>
              <p className="text-xs text-gray-500 uppercase font-semibold">
                Total Quantity
              </p>
              <p className="text-gray-800 font-medium">
                {orders.total_quantity} pcs
              </p>
            </div>
            <div>
              <p className="text-xs text-gray-500 uppercase font-semibold">
                Current Factory Location
              </p>
              <div className="font-semibold text-gray-700">
                {getCurrentLocation(orders.status, orders.layout_id)}
              </div>
            </div>
          </div>

          {/* Timelines */}
          <div className="bg-white p-6 rounded-lg shadow space-y-4">
            <h2 className="text-lg font-bold text-gray-800 border-b pb-2">
              Sampling & Timelines
            </h2>
            <div>
              <p className="text-xs text-gray-500 uppercase font-semibold">
                Internal Deadline (4 Days)
              </p>
              <p className="text-gray-800 font-medium">
                {new Date(orders.internal_sample_deadline).toLocaleDateString()}
              </p>
            </div>
            <div>
              <p className="text-xs text-gray-500 uppercase font-semibold">
                Customer Deadline (7 Days)
              </p>
              <p className="text-gray-800 font-medium">
                {new Date(orders.customer_sample_deadline).toLocaleDateString()}
              </p>
            </div>
            <div>
              <p className="text-xs text-gray-500 uppercase font-semibold">
                Actual Sample Finished At
              </p>
              <p className="font-medium text-gray-800">
                {orders.actual_sample_finished_at ? (
                  new Date(orders.actual_sample_finished_at).toLocaleString()
                ) : (
                  <span className="text-yellow-600">Pending completion</span>
                )}
              </p>
            </div>
          </div>
          <div className="bg-white p-6 rounded-lg shadow space-y-4">
            <h2 className="text-lg font-bold text-gray-800 border-b pb-2">
              Mass Production Timelines
            </h2>
            <div className="grid grid-cols-2 gap-4 text-sm">
              <div>
                <p className="text-gray-500 font-semibold">
                  INTERNAL TARGET (9 DAYS)
                </p>
                <p className="text-gray-900">
                  {orders.internal_production_deadline
                    ? new Date(
                        orders.internal_production_deadline,
                      ).toLocaleDateString()
                    : "Not started"}
                </p>
              </div>
              <div>
                <p className="text-gray-500 font-semibold">
                  CUSTOMER TARGET (12 DAYS)
                </p>
                <p className="text-gray-900">
                  {orders.customer_production_deadline
                    ? new Date(
                        orders.customer_production_deadline,
                      ).toLocaleDateString()
                    : "Not started"}
                </p>
              </div>
              <div>
                <p className="text-gray-500 font-semibold">
                  ACTUAL PRODUCTION START
                </p>
                <p className="text-gray-900">
                  {orders.actual_production_started_at
                    ? new Date(
                        orders.actual_production_started_at,
                      ).toLocaleString()
                    : "Pending..."}
                </p>
              </div>
              <div>
                <p className="text-gray-500 font-semibold">
                  ACTUAL COMPLETION (SHIPPING)
                </p>
                <p className="text-gray-900">
                  {orders.actual_production_completed_at
                    ? new Date(
                        orders.actual_production_completed_at,
                      ).toLocaleString()
                    : "In Progress..."}
                </p>
              </div>
            </div>
          </div>
        </div>

        {/* Size Breakdown Card */}
        <div className="bg-white p-6 rounded-lg shadow">
          <h2 className="text-lg font-bold text-gray-800 border-b pb-2 mb-4">
            Size Breakdown
          </h2>
          <div className="flex flex-wrap gap-3">
            {orders.sizes && orders.sizes.length > 0 ? (
              orders.sizes.map((s, idx) => (
                <div
                  key={idx}
                  className="bg-gray-100 border border-gray-200 rounded-lg p-3 text-center min-w-22.5"
                >
                  <p className="text-xs text-gray-500 font-bold uppercase">
                    Size {s.size_label}
                  </p>
                  <p className="text-lg font-bold text-gray-800">
                    {s.quantity}
                  </p>
                </div>
              ))
            ) : (
              <p className="text-gray-500 text-sm">
                No size quantities recorded.
              </p>
            )}
          </div>
        </div>

        {/* Tech Packs / Images Gallery */}
        <div className="bg-white p-6 rounded-lg shadow">
          <h2 className="text-lg font-bold text-gray-800 border-b pb-2 mb-4">
            Tech Pack Images & Sketches
          </h2>
          <div className="grid grid-cols-2 sm:grid-cols-3 gap-4">
            {orders.images && orders.images.length > 0 ? (
              orders.images.map((img) => (
                <div
                  key={img.id}
                  className="border rounded-lg overflow-hidden bg-gray-50"
                >
                  {/* Pointing directly to Go static server uploads */}
                  <Image
                    src={`http://localhost:8080/${img.image_url}`}
                    alt="Tech Pack"
                    className="w-full h-40 object-cover"
                  />
                  <div className="p-2 text-xs text-gray-500 text-center truncate">
                    Image ID: #{img.id}
                  </div>
                </div>
              ))
            ) : (
              <p className="text-gray-500 text-sm col-span-full">
                No tech pack images uploaded for this order yet.
              </p>
            )}
          </div>
        </div>
      </div>
    </main>
  );
}
