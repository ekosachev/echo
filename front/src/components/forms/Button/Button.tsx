import { cn } from '@/lib/utils'
import * as React from 'react'

export interface ButtonProps extends React.HTMLAttributes<HTMLButtonElement> {
  type?: 'submit' | 'reset' | 'button' | undefined,
  variant?: 'primary' | 'secondary',
  disabled?: boolean,
}

const Button = React.forwardRef<HTMLButtonElement, ButtonProps>((
  {
    className,
    type = undefined,
    variant = 'primary',
    disabled = false,
    ...props
  }, ref
) => {
  return (
    <button
      type={type}
      ref={ref}
      className={cn(
        variant === 'primary' && 'bg-sky-800 text-sky-50 hover:bg-sky-500',
        variant === 'secondary' && 'border-2 bg-transparent text-sky-800 border-sky-800 hover:border-sky-700 hover:bg-sky-700 hover:text-sky-50',
        'focus:outline-none',
        'transition-colors duration-75',
        'rounded-lg',
        'px-6 py-2',
        'text-lg',
        className,
      )}
      {...props}
    />
  )
})

Button.displayName = 'Button'

export { Button }
