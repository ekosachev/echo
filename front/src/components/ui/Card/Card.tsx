import { cn } from '@/lib/utils'
import * as React from 'react'

// 1. Определяем базовые типы пропсов
export interface CardProps extends React.HTMLAttributes<HTMLDivElement> {
  variant?: 'default' | 'outline' | 'ghost'
  hoverable?: boolean
  padding?: 'none' | 'sm' | 'md' | 'lg'
  radius?: 'none' | 'sm' | 'md' | 'lg' | 'full'
  shadow?: 'none' | 'sm' | 'md' | 'lg'
  bordered?: boolean
  disabled?: boolean
}

// 2. Создаем корневой компонент Card
const Card = React.forwardRef<HTMLDivElement, CardProps>(
  ({
    className,
    variant = 'default',
    hoverable = false,
    padding = 'md',
    radius = 'md',
    shadow = 'sm',
    bordered = true,
    disabled = false,
    ...props
  }, ref) => {
    // 3. Динамические классы на основе пропсов
    const variantClasses = {
      default: 'dark:bg-slate-900/50',
      outline: 'bg-transparent border',
      ghost: 'bg-transparent shadow-none'
    }

    const paddingClasses = {
      none: 'p-0',
      sm: 'p-3',
      md: 'p-4 md:p-5',
      lg: 'p-6 md:p-8'
    }

    const radiusClasses = {
      none: 'rounded-none',
      sm: 'rounded-sm',
      md: 'rounded-lg',
      lg: 'rounded-xl',
      full: 'rounded-2xl'
    }

    const shadowClasses = {
      none: 'shadow-none',
      sm: 'shadow-sm',
      md: 'shadow',
      lg: 'shadow-md'
    }

    return (
      <div
        ref={ref}
        className={cn(
          // Базовые стили
          variantClasses[variant],
          paddingClasses[padding],
          radiusClasses[radius],
          shadowClasses[shadow],
          bordered && 'border border-slate-400 dark:border-slate-800',
          hoverable && !disabled && 'hover:shadow-lg hover:border-sky-300 cursor-pointer hover:transition-all duration-75',
          disabled && 'opacity-50 cursor-not-allowed',
          className // Позволяет переопределить стили снаружи
        )}
        aria-disabled={disabled}
        {...props}
      />
    )
  }
)

Card.displayName = 'Card'

export { Card }
