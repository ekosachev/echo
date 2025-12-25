import { cn } from '@/lib/utils'
import * as React from 'react'

export interface FieldProps extends React.HTMLAttributes<HTMLDivElement> {
  name: string
  type: 'password' | 'text' | 'email'
  disabled?: boolean
  label?: string
  placeholder?: string
  autoComplete?: string
  required?: boolean
}

const Field = React.forwardRef<HTMLInputElement, FieldProps>((
  {
    className,
    name,
    type,
    disabled = false,
    label = '',
    placeholder = '',
    autoComplete = '',
    required = false,
    ...props
  }, ref
) => {
  return (
    <div ref={ref} className={cn(className, 'flex flex-col p-2 gap-1',)} >
      <input type={type} id={'id' + name} name={name} placeholder={placeholder} required={required} {...props} className='peer text-lg border-2 rounded-md border-slate-700 focus:border-sky-300 transaition-all duration-75 focus:outline-none p-1' autoComplete={autoComplete}></input>
      {label !== '' && (<label htmlFor={'id' + name} className='text-md text-slate-700 peer-focus:text-sky-300'>{label}</label>)}

    </div>
  )
})

Field.displayName = 'Field'

export { Field }

