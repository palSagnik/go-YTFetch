import * as React from "react";

const buttonVariants = {
  variant: {
    default:
      "bg-blue-500 text-white hover:bg-blue-600 hover:shadow-lg transform hover:-translate-y-0.5 transition-all duration-300",
    destructive:
      "bg-red-500 text-white hover:bg-red-600 hover:shadow-lg transform hover:-translate-y-0.5 transition-all duration-300",
    outline:
      "border border-blue-500 text-blue-500 hover:bg-blue-50 hover:text-blue-600 hover:border-blue-600 transition-all duration-300",
    secondary:
      "bg-gray-200 text-gray-800 hover:bg-gray-300 hover:shadow-md transform hover:-translate-y-0.5 transition-all duration-300",
    ghost:
      "text-gray-600 hover:bg-gray-100 hover:text-gray-800 transition-all duration-300",
    link: "text-blue-500 underline-offset-4 hover:underline transition-all duration-300",
  },
  size: {
    default: "h-10 px-4 py-2",
    sm: "h-9 rounded-md px-3",
    lg: "h-11 rounded-md px-8",
    icon: "h-10 w-10",
  },
};

export interface ButtonProps
  extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: keyof typeof buttonVariants.variant;
  size?: keyof typeof buttonVariants.size;
}

const Button = React.forwardRef<HTMLButtonElement, ButtonProps>(
  ({ className, variant = "default", size = "default", ...props }, ref) => {
    const baseClasses =
      "inline-flex items-center justify-center rounded-md text-sm font-medium ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50";
    const variantClasses = buttonVariants.variant[variant];
    const sizeClasses = buttonVariants.size[size];
    const allClasses = `${baseClasses} ${variantClasses} ${sizeClasses} ${
      className || ""
    }`;

    return <button className={allClasses} ref={ref} {...props} />;
  }
);
Button.displayName = "Button";

export { Button, buttonVariants };
