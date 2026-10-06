import OrderForm from "@/components/OrderForm";
import Link from "next/link";

export default function CreateOrderPage() {
  return (
    <div className="max-w-7xl py-6 mx-auto">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-3xl font-extrabold text-gray-50">Create New Order</h1>
        <Link href="/orders" className="text-gray-500 hover:underline font-medium">
          Cancel
        </Link>
      </div>

      <OrderForm />
    </div>
  )
}