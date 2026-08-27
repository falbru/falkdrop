interface Props {
  children: string;
}

const Title = (props: Props) => {
  const { children } = props;

  return <title>{`${children} — FalkDrop`}</title>;
};

export default Title;
