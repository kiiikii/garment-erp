import Link from "next/link";

// Define the shape of our order data
interface Order {
  id: number;
  status: string;
  total_quantity: number;
  waiting_reason?: string | null;
  waiting_started_at?: string | null;
  internal_production_deadline?: string | null;
  customer_production_deadline?: string | null;
  actual_production_completed_at?: string | null;
}

export default async function Home() {
  let orders: Order[] = [];
  let fetchError = false;

  try {
    // Fetch all orders on the server (no-store ensures live data)
    const res = await fetch("http://localhost:8080/orders", {
      cache: "no-store",
    });
    if (res.ok) {
      orders = await res.json();
    } else {
      fetchError = true;
    }
  } catch (err) {
    fetchError = true;
  }

  // ==========================================
  // SPRINT 15: MONITORING CALCULATIONS
  // ==========================================

  // 1. Filter out completed or rejected orders
  const activeOrders = orders.filter(
    (o) => o.status !== "SHIPPED" && !o.status.includes("REJECTED"),
  );

  // 15.9 Bottlenecks: Orders currently halted on the floor
  const bottlenecks = activeOrders.filter((o) => o.waiting_reason);

  // 15.5 Delays & 15.3 Variance: Orders past their internal deadline
  const today = new Date();
  const delayedOrders = activeOrders.filter((o) => {
    if (!o.internal_production_deadline) return false;
    const deadline = new Date(o.internal_production_deadline);
    // If today is past the deadline and it's not finished
    return today > deadline && !o.actual_production_completed_at;
  });

  // 15.8 Current Process: Group orders by their exact location
  const processCounts = activeOrders.reduce(
    (acc: Record<string, number>, o) => {
      acc[o.status] = (acc[o.status] || 0) + 1;
      return acc;
    },
    {},
  );

  // 15.4 Completion Volume: Total garments currently in progress
  const totalWipGarments = activeOrders.reduce(
    (sum, o) => sum + o.total_quantity,
    0,
  );

  const getProgress = (status: string) => {
    const progressMap: Record<string, number> = {
      WAITING_FOR_SAMPLE: 5, SAMPLE_REVISED: 5, SAMPLE_APPROVED: 10,
      LOA_REVISED: 15, LOA_SIGNED: 20, WAITING_FOR_MATERIALS: 20,
      READY_FOR_PRODUCTION: 25, IN_PRODUCTION: 30, CUTTING: 45,
      SEWING: 65, QC: 85, FINISHING: 95, READY_FOR_SHIPPING: 98, SHIPPED: 100
    };
    return progressMap[status] || 0;
  };

  return (
    <main className="min-h-screen bg-gray-100 p-8 font-sans">
      <div className="max-w-6xl mx-auto space-y-8">
        {/* Header & Quick Links */}
        <div className="flex justify-between items-center bg-white p-6 rounded-lg shadow-sm">
          <div>
            <h1 className="text-3xl font-extrabold text-gray-900">
              Factory Command Center
            </h1>
          </div>
          <div className="flex gap-4">
            <Link
              href="/customers"
              className="bg-white border-2 border-gray-300 text-gray-700 px-6 py-2 rounded font-bold hover:bg-gray-50 transition"
            >
              Customers Directory
            </Link>
            <Link
              href="/orders"
              className="bg-blue-600 text-white px-6 py-2 rounded font-bold hover:bg-blue-700 shadow transition"
            >
              All Orders
            </Link>
          </div>
        </div>

        {fetchError && (
          <div className="bg-red-100 text-red-800 p-4 rounded-lg font-bold">
            Failed to connect to Go backend. Make sure localhost:8080 is
            running.
          </div>
        )}

        {/* SPRINT 15: KPI CARDS */}
        <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
          <div className="bg-white p-6 rounded-lg shadow-sm border-l-4 border-blue-500">
            <h3 className="text-sm font-bold text-gray-500 uppercase">
              Active Orders
            </h3>
            <p className="text-4xl font-extrabold text-gray-900 mt-2">
              {activeOrders.length}
            </p>
          </div>

          <div className="bg-white p-6 rounded-lg shadow-sm border-l-4 border-indigo-500">
            <h3 className="text-sm font-bold text-gray-500 uppercase">
              WIP Garments
            </h3>
            <p className="text-4xl font-extrabold text-gray-900 mt-2">
              {totalWipGarments}{" "}
              <span className="text-lg text-gray-400 font-medium">pcs</span>
            </p>
          </div>

          <div
            className={`bg-white p-6 rounded-lg shadow-sm border-l-4 ${bottlenecks.length > 0 ? "border-red-500" : "border-green-500"}`}
          >
            <h3 className="text-sm font-bold text-gray-500 uppercase">
              Bottlenecks (Halted)
            </h3>
            <p
              className={`text-4xl font-extrabold mt-2 ${bottlenecks.length > 0 ? "text-red-600" : "text-green-600"}`}
            >
              {bottlenecks.length}
            </p>
          </div>

          <div
            className={`bg-white p-6 rounded-lg shadow-sm border-l-4 ${delayedOrders.length > 0 ? "border-orange-500" : "border-green-500"}`}
          >
            <h3 className="text-sm font-bold text-gray-500 uppercase">
              Delayed Orders
            </h3>
            <p
              className={`text-4xl font-extrabold mt-2 ${delayedOrders.length > 0 ? "text-orange-600" : "text-green-600"}`}
            >
              {delayedOrders.length}
            </p>
          </div>
        </div>

        {/* SPRINT 15: PIPELINE VISIBILITY */}
        <div className="bg-white p-6 rounded-lg shadow-sm">
          <h2 className="text-xl font-bold text-gray-800 mb-6">
            Current Process Distribution
          </h2>
          <div className="grid grid-cols-2 md:grid-cols-5 gap-4">
            {/* Sampling Phase */}
            <div className="bg-gray-50 p-4 rounded border text-center">
              <span className="block text-xs font-bold text-gray-500 mb-1">
                SAMPLING
              </span>
              <span className="text-2xl font-bold text-gray-800">
                {(processCounts["WAITING_FOR_SAMPLE"] || 0) +
                  (processCounts["SAMPLE_REVISED"] || 0)}
              </span>
            </div>

            {/* Waiting for Materials */}
            <div className="bg-gray-50 p-4 rounded border text-center">
              <span className="block text-xs font-bold text-gray-500 mb-1">
                MATERIAL WAIT
              </span>
              <span className="text-2xl font-bold text-gray-800">
                {processCounts["WAITING_FOR_MATERIALS"] || 0}
              </span>
            </div>

            {/* Cutting */}
            <div className="bg-blue-50 p-4 rounded border border-blue-200 text-center">
              <span className="block text-xs font-bold text-blue-600 mb-1">
                CUTTING
              </span>
              <span className="text-2xl font-bold text-blue-900">
                {processCounts["CUTTING"] || 0}
              </span>
            </div>

            {/* Sewing */}
            <div className="bg-blue-50 p-4 rounded border border-blue-200 text-center">
              <span className="block text-xs font-bold text-blue-600 mb-1">
                SEWING
              </span>
              <span className="text-2xl font-bold text-blue-900">
                {processCounts["SEWING"] || 0}
              </span>
            </div>

            {/* QC & Finishing */}
            <div className="bg-green-50 p-4 rounded border border-green-200 text-center">
              <span className="block text-xs font-bold text-green-600 mb-1">
                QC & FINISH
              </span>
              <span className="text-2xl font-bold text-green-900">
                {(processCounts["QC"] || 0) + (processCounts["FINISHING"] || 0)}
              </span>
            </div>
          </div>
        </div>

        {/* BOTTLE NECK ALERT LIST */}
        {bottlenecks.length > 0 && (
          <div className="bg-red-50 p-6 rounded-lg shadow-sm border border-red-200">
            <h2 className="text-xl font-bold text-red-800 mb-4 flex items-center gap-2">
              Action Required: Halted Production
            </h2>
            <div className="space-y-4">
              {bottlenecks.map(order => {
                const progress = getProgress(order.status);
                
                // 15.6 Waiting Time Calculation
                let daysWaiting = 0;
                if (order.waiting_started_at) {
                  const waitStart = new Date(order.waiting_started_at);
                  const diffTime = Math.abs(today.getTime() - waitStart.getTime());
                  daysWaiting = Math.ceil(diffTime / (1000 * 60 * 60 * 24));
                }

                // 15.7 Deadline Formatting
                const custDeadline = order.customer_production_deadline 
                  ? new Date(order.customer_production_deadline).toLocaleDateString() 
                  : "N/A";

                return (
                  <div key={order.id} className="bg-white p-5 rounded border border-red-200 shadow-sm">
                    <div className="flex justify-between items-start mb-3">
                      <div>
                        <span className="font-extrabold text-gray-900 text-lg">Order #{order.id}</span>
                        <span className="bg-red-100 text-red-800 text-xs font-bold px-2 py-1 rounded ml-3">
                          Stuck in {order.status}
                        </span>
                        <p className="text-red-600 font-bold mt-2">🚨 {order.waiting_reason}</p>
                        
                        {/* 15.6 & 15.7 Rendered Here */}
                        <div className="text-sm text-gray-500 mt-2 flex gap-4">
                          <span className="font-medium text-red-500">
                            ⏳ Waiting Time: {daysWaiting} days
                          </span>
                          <span>|</span>
                          <span>Customer Deadline: <span className="font-bold text-gray-700">{custDeadline}</span></span>
                        </div>
                      </div>
                      <Link href={`/orders/${order.id}`} className="bg-red-600 text-white px-4 py-2 rounded font-bold text-sm hover:bg-red-700">
                        Resolve Issue &rarr;
                      </Link>
                    </div>

                    {/* 15.4 Completion Percentage Bar */}
                    <div className="mt-4">
                      <div className="flex justify-between text-xs mb-1 font-bold text-gray-500">
                        <span>Est. Completion</span>
                        <span>{progress}%</span>
                      </div>
                      <div className="w-full bg-gray-200 rounded-full h-2.5">
                        <div className="bg-blue-600 h-2.5 rounded-full transition-all" style={{ width: `${progress}%` }}></div>
                      </div>
                    </div>
                  </div>
                );
              })}
            </div>
          </div>
        )}
      </div>
    </main>
  );
}
